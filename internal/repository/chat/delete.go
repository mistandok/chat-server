package chat

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	serviceModel "github.com/mistandok/chat-server/internal/model"
)

// Delete ..
func (r *Repo) Delete(ctx context.Context, chatID serviceModel.ChatID) error {
	queryFormat := `DELETE FROM %s WHERE %s = @%s`
	query := fmt.Sprintf(queryFormat, chatTable, idColumn, idColumn)

	args := pgx.NamedArgs{
		idColumn: chatID,
	}

	_, err := r.pool.Exec(ctx, query, args)

	return err
}
