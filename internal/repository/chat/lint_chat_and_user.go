package chat

import (
	"context"

	"github.com/jackc/pgx/v5"
	serviceModel "github.com/mistandok/chat-server/internal/model"
	"github.com/mistandok/chat-server/internal/repository/chat/convert"
)

// LinkChatAndUsers ..
func (r *Repo) LinkChatAndUsers(ctx context.Context, chatID serviceModel.ChatID, userIDs []serviceModel.UserID) error {
	repoChatID := int64(chatID)
	repoUserIDs := convert.ToSliceIntFromSliceServiceUserID(userIDs)

	countUsers := len(repoUserIDs)
	if countUsers == 0 {
		return nil
	}

	rows := make([][]interface{}, 0)
	for _, userID := range repoUserIDs {
		rows = append(rows, []interface{}{repoChatID, userID})
	}

	_, err := r.pool.CopyFrom(
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
