package user

import (
	"context"
	"fmt"
	"strings"

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
	_, err := r.pool.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
