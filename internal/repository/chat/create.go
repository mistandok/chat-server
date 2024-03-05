package chat

import (
	"context"
	"fmt"

	"github.com/mistandok/chat-server/internal/client/db"

	"github.com/jackc/pgx/v5"
)

// Create chat in db
func (r *Repo) Create(ctx context.Context) (int64, error) {
	queryFormat := `INSERT INTO %s DEFAULT VALUES RETURNING id`
	query := fmt.Sprintf(queryFormat, chatTable)
	q := db.Query{
		Name:     "chat_repository.Create",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	chatID, err := pgx.CollectOneRow(rows, pgx.RowTo[int64])
	if err != nil {
		return 0, err
	}

	return chatID, nil
}
