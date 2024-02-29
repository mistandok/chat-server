package message

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	serviceModel "github.com/mistandok/chat-server/internal/model"
	"github.com/mistandok/chat-server/internal/repository"
	"github.com/mistandok/chat-server/internal/repository/message/convert"
	"github.com/pkg/errors"
)

// Create ..
func (r *Repo) Create(ctx context.Context, message serviceModel.Message) error {
	repoMessage := convert.ToMessageFromServiceMessage(&message)

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

	_, err := r.pool.Exec(ctx, query, args)
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
