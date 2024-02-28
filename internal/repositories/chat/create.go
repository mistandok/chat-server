package chat

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"strings"
)

// Create chat in db
func (c *Repo) Create(ctx context.Context, userIDs []int64) (int64, error) {
	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			err = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	chatID, err := c.createChatForUsers(ctx, tx, userIDs)
	if err != nil {
		return 0, errors.Errorf("ошибка при генерации пользовательского чата: %v", err)
	}

	return chatID, nil
}

func (c *Repo) createChatForUsers(ctx context.Context, tx pgx.Tx, userIDs []int64) (int64, error) {
	chatID, err := c.createChat(ctx, tx)
	if err != nil {
		return 0, fmt.Errorf("ошибка при создании чата: %v", err)
	}

	err = c.createUsers(ctx, tx, userIDs)
	if err != nil {
		return 0, fmt.Errorf("ошибка при создании пользователей: %v", err)
	}

	err = c.linkChatAndUsers(ctx, tx, chatID, userIDs)
	if err != nil {
		return 0, fmt.Errorf("ошибка при связывании чата с пользователями: %v", err)
	}

	return chatID, nil
}

func (c *Repo) createChat(ctx context.Context, tx pgx.Tx) (int64, error) {
	queryFormat := `INSERT INTO %s DEFAULT VALUES RETURNING id`
	query := fmt.Sprintf(queryFormat, chatTable)

	rows, err := tx.Query(ctx, query)
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

func (c *Repo) createUsers(ctx context.Context, tx pgx.Tx, userIDs []int64) error {
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
	_, err := tx.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (c *Repo) linkChatAndUsers(ctx context.Context, tx pgx.Tx, chatID int64, userIDs []int64) error {
	countUsers := len(userIDs)
	if countUsers == 0 {
		return nil
	}

	rows := make([][]interface{}, 0)
	for _, userID := range userIDs {
		rows = append(rows, []interface{}{chatID, userID})
	}

	_, err := tx.CopyFrom(
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
