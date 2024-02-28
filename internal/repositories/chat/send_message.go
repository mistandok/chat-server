package chat

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mistandok/chat-server/internal/repositories"
	"github.com/mistandok/chat-server/internal/repositories/chat/model"
	"github.com/pkg/errors"
)

// SendMessage save message in db
func (c *Repo) SendMessage(ctx context.Context, in *model.Message) error {
	userLinkedWithChat, err := c.isUserInChat(ctx, in.ToChatID, in.FromUserID)
	if err != nil {
		return errors.Errorf("ошибка во время проверки наличия пользователя в чате: %v", err)
	}

	if !userLinkedWithChat {
		return repositories.ErrUserNotInTheChat
	}

	queryFormat := `
		INSERT INTO %s (%s, %s, %s, %s)
		VALUES (@%s, @%s, @%s, @%s)
    `
	query := fmt.Sprintf(
		queryFormat,
		messageTable,
		fromUserIDColumn, chatIDColumn, messageColumn, sentAtColumn,
		fromUserIDColumn, chatIDColumn, messageColumn, sentAtColumn,
	)

	args := pgx.NamedArgs{
		fromUserIDColumn: in.FromUserID,
		chatIDColumn:     in.ToChatID,
		messageColumn:    in.Message,
		sentAtColumn:     in.SendTime,
	}

	_, err = c.pool.Exec(ctx, query, args)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch {
			case pgErr.ConstraintName == messageChatIDFKConstraint:
				return repositories.ErrChatNotFound
			case pgErr.ConstraintName == messageFromUserIDFKConstraint:
				return repositories.ErrUserNotFound
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

	c.logger.Info().Msg(query)

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
