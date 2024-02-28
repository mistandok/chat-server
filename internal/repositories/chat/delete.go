package chat

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

func (c *Repo) Delete(ctx context.Context, chatID int64) error {
	queryFormat := `DELETE FROM %s WHERE %s = @%s`
	query := fmt.Sprintf(queryFormat, chatTable, idColumn, idColumn)

	args := pgx.NamedArgs{
		idColumn: chatID,
	}

	_, err := c.pool.Exec(ctx, query, args)

	return err
}
