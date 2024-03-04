package chat

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// LinkChatAndUsers ..
func (r *Repo) LinkChatAndUsers(ctx context.Context, chatID int64, userIDs []int64) error {
	countUsers := len(userIDs)
	if countUsers == 0 {
		return nil
	}

	rows := make([][]interface{}, 0)
	for _, userID := range userIDs {
		rows = append(rows, []interface{}{chatID, userID})
	}

	_, err := r.db.DB().CopyFromContext(
		ctx,
		pgx.Identifier{chatUserTable},
		[]string{chatIDColumn, userIDColumn},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return err
	}

	return nil
}
