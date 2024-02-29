package chat

import (
	"context"
	"fmt"

	"github.com/mistandok/chat-server/internal/client/db"

	"github.com/jackc/pgx/v5"
	serviceModel "github.com/mistandok/chat-server/internal/model"
)

// Delete ..
func (r *Repo) Delete(ctx context.Context, chatID serviceModel.ChatID) error {
	queryFormat := `DELETE FROM %s WHERE %s = @%s`
	query := fmt.Sprintf(queryFormat, chatTable, idColumn, idColumn)
	q := db.Query{
		Name:     "chat_repository.Delete",
		QueryRaw: query,
	}
	args := pgx.NamedArgs{
		idColumn: chatID,
	}

	_, err := r.db.DB().ExecContext(ctx, q, args)

	return err
}
