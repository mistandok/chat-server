package chat

import (
	"context"
	"errors"
	"fmt"

	"github.com/mistandok/chat-server/internal/client/db"

	"github.com/jackc/pgx/v5"
)

// IsUserInChat ..
func (r *Repo) IsUserInChat(ctx context.Context, chatID int64, userID int64) (bool, error) {
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

	q := db.Query{
		Name:     "chat_repository.IsUserInChat",
		QueryRaw: query,
	}

	args := pgx.NamedArgs{
		chatIDColumn: chatID,
		userIDColumn: userID,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args)
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
