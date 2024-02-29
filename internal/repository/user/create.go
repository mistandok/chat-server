package user

import (
	"context"
	"fmt"
	"strings"

	"github.com/mistandok/chat-server/internal/client/db"

	serviceModel "github.com/mistandok/chat-server/internal/model"
	"github.com/mistandok/chat-server/internal/repository/user/convert"
)

// CreateMass ..
func (r *Repo) CreateMass(ctx context.Context, userIDs []serviceModel.UserID) error {
	repoUserIDs := convert.ToSliceIntFromSliceServiceUserID(userIDs)

	countUsers := len(repoUserIDs)
	if countUsers == 0 {
		return nil
	}

	var strUserIDs = make([]string, 0, countUsers)
	for _, userID := range repoUserIDs {
		strUserIDs = append(strUserIDs, fmt.Sprintf("(%d)", userID))
	}

	values := strings.Join(strUserIDs, ",")

	queryFormat := `
		INSERT INTO "%s" (%s)
		VALUES %s
		ON CONFLICT DO NOTHING
	`
	query := fmt.Sprintf(queryFormat, userTable, idColumn, values)
	q := db.Query{
		Name:     "user_repository.CreateMass",
		QueryRaw: query,
	}

	_, err := r.db.DB().ExecContext(ctx, q)
	if err != nil {
		return err
	}

	return nil
}
