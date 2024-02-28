package chat

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	serviceModel "github.com/mistandok/chat-server/internal/model"
	"github.com/mistandok/chat-server/internal/repository"
	"github.com/mistandok/chat-server/internal/repository/chat/convert"
	"github.com/pkg/errors"
)

// SendMessage save message in db
func (c *Repo) SendMessage(ctx context.Context, message serviceModel.Message) error {
	repoMessage := convert.ToMessageFromServiceMessage(&message)

	userLinkedWithChat, err := c.isUserInChat(ctx, repoMessage.ToChatID, repoMessage.FromUserID)
	if err != nil {
		return fmt.Errorf("ошибка во время проверки наличия пользователя в чате: %w", err)
	}

	if !userLinkedWithChat {
		return repository.ErrUserNotInTheChat
	}

	queryFormat := `
		INSERT INTO %s (%s, %s, %s, %s)
		VALUES (@%s, @%s, @%s, @%s)
    `
	query := fmt.Sprintf(
		queryFormat,
		messageTable,
		fromUserIDColumn, chatIDColumn, textColumn, sentAtColumn,
		fromUserIDColumn, chatIDColumn, textColumn, sentAtColumn,
	)

	args := pgx.NamedArgs{
		fromUserIDColumn: repoMessage.FromUserID,
		chatIDColumn:     repoMessage.ToChatID,
		textColumn:       repoMessage.Text,
		sentAtColumn:     repoMessage.SendTime,
	}

	_, err = c.pool.Exec(ctx, query, args)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch {
			case pgErr.ConstraintName == messageChatIDFKConstraint:
				return repository.ErrChatNotFound
			case pgErr.ConstraintName == messageFromUserIDFKConstraint:
				return repository.ErrUserNotFound
			}
		}
		return err
	}

	return nil
}

func (c *Repo) isUserInChat(ctx context.Context, chatID int64, userID int64) (bool, error) {
	queryFormat := `
		SELECT TRUE
		FROM %s
		WHERE %s = @%s AND %s = @%s
    `
	query := fmt.Sprintf(
		queryFormat,
		chatUserTable,
		chatIDColumn, chatIDColumn, userIDColumn, userIDColumn,
	)

	args := pgx.NamedArgs{
		chatIDColumn: chatID,
		userIDColumn: userID,
	}

	rows, err := c.pool.Query(ctx, query, args)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	_, err = pgx.CollectOneRow(rows, pgx.RowTo[bool])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
