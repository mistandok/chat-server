package chat

import (
	"context"
	"fmt"

	"github.com/mistandok/chat-server/internal/client/db"

	"github.com/jackc/pgx/v5"
	serviceModel "github.com/mistandok/chat-server/internal/model"
	"github.com/pkg/errors"
)

// IsUserInChat ..
func (r *Repo) IsUserInChat(ctx context.Context, chatID serviceModel.ChatID, userID serviceModel.UserID) (bool, error) {
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
		chatIDColumn: int64(chatID),
		userIDColumn: int64(userID),
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
