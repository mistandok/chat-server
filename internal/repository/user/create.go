package user

import (
	"context"
	"fmt"
	"strings"

	"github.com/mistandok/platform_common/pkg/db"
)

// CreateMass ..
func (r *Repo) CreateMass(ctx context.Context, userIDs []int64) error {
	countUsers := len(userIDs)
	if countUsers == 0 {
		return nil
	}

	var strUserIDs = make([]string, 0, countUsers)
	for _, userID := range userIDs {
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
