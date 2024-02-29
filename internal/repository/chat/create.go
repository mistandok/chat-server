package chat

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	serviceModel "github.com/mistandok/chat-server/internal/model"
)

// Create chat in db
func (r *Repo) Create(ctx context.Context) (serviceModel.ChatID, error) {
	queryFormat := `INSERT INTO %s DEFAULT VALUES RETURNING id`
	query := fmt.Sprintf(queryFormat, chatTable)

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	chatID, err := pgx.CollectOneRow(rows, pgx.RowTo[int64])
	if err != nil {
		return 0, err
	}

	return serviceModel.ChatID(chatID), nil
}
