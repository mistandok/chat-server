package postgresql

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mistandok/chat-server/internal/repositories"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// ChatRepo user repo for operations with chat.
type ChatRepo struct {
	pool   *pgxpool.Pool
	logger *zerolog.Logger
}

// NewChatRepo  get new repo instance.
func NewChatRepo(pool *pgxpool.Pool, logger *zerolog.Logger) *ChatRepo {
	return &ChatRepo{
		pool:   pool,
		logger: logger,
	}
}

// Create chat in db
func (c *ChatRepo) Create(ctx context.Context, in *repositories.ChatCreateIn) (*repositories.ChatCreateOut, error) {
	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	chatID, err := c.createChatForUsers(ctx, tx, in.UserIDs)
	if err != nil {
		return nil, errors.Errorf("ошибка при генерации пользовательского чата: %v", err)
	}

	return &repositories.ChatCreateOut{ID: chatID}, nil
}

// Delete delete chat from db.
func (c *ChatRepo) Delete(ctx context.Context, in *repositories.ChatDeleteIn) error {
	query := `DELETE FROM chat WHERE id = @id`

	args := pgx.NamedArgs{
		"id": in.ID,
	}

	_, err := c.pool.Exec(ctx, query, args)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// SendMessage save message in db
func (c *ChatRepo) SendMessage(ctx context.Context, in *repositories.SendMessageIn) error {
	userLinkedWithChat, err := c.isUserInChat(ctx, in.ToChatID, in.FromUserID)
	if err != nil {
		return errors.Errorf("ошибка во время проверки наличия пользователя в чате: %v", err)
	}

	if !userLinkedWithChat {
		return repositories.ErrUserNotInTheChat
	}

	query := `
		INSERT INTO message (from_user_id, chat_id, message, sended)
		VALUES (@fromUserID, @toChatID, @message, @sended)
    `

	args := pgx.NamedArgs{
		"fromUserID": in.FromUserID,
		"toChatID":   in.ToChatID,
		"message":    in.Message,
		"sended":     in.SendTime,
	}

	_, err = c.pool.Exec(ctx, query, args)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch {
			case pgErr.ConstraintName == "fk_chat_id":
				return repositories.ErrChatNotFound
			case pgErr.ConstraintName == "fk_from_user_id":
				return repositories.ErrUserNotFound
			}
		}
		return errors.WithStack(err)
	}

	return nil
}

func (c *ChatRepo) createChatForUsers(ctx context.Context, tx pgx.Tx, userIDs []int64) (int64, error) {
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

func (c *ChatRepo) createChat(ctx context.Context, tx pgx.Tx) (int64, error) {
	query := `INSERT INTO chat DEFAULT VALUES RETURNING id`

	rows, err := tx.Query(ctx, query)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	defer rows.Close()

	chat, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[repositories.ChatCreateOut])
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return chat.ID, nil
}

func (c *ChatRepo) createUsers(ctx context.Context, tx pgx.Tx, userIDs []int64) error {
	countUsers := len(userIDs)
	if countUsers == 0 {
		return nil
	}

	tempTableName, err := c.createAndFillTempUserTable(ctx, tx, userIDs)
	if err != nil {
		return errors.Errorf("ошибка во время создания временной таблицы пользователей: %v", err)
	}

	query := fmt.Sprintf(`INSERT INTO "user" SELECT * FROM %s ON CONFLICT DO NOTHING`, tempTableName)

	_, err = tx.Exec(ctx, query)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *ChatRepo) createAndFillTempUserTable(ctx context.Context, tx pgx.Tx, userIDs []int64) (string, error) {
	tableName := "_temp_upsert_user"
	query := fmt.Sprintf(`CREATE TEMPORARY TABLE %s (LIKE "user" INCLUDING ALL) ON COMMIT DROP`, tableName)

	_, err := tx.Exec(ctx, query)
	if err != nil {
		return "", errors.WithStack(err)
	}

	rows := make([][]interface{}, 0, len(userIDs))
	for _, userID := range userIDs {
		rows = append(rows, []interface{}{userID})
	}

	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{tableName},
		[]string{"id"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return "", err
	}

	return tableName, nil
}

func (c *ChatRepo) linkChatAndUsers(ctx context.Context, tx pgx.Tx, chatID int64, userIDs []int64) error {
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
		pgx.Identifier{"chat_user"},
		[]string{"chat_id", "user_id"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *ChatRepo) isUserInChat(ctx context.Context, chatID int64, userID int64) (bool, error) {
	query := `
		SELECT
			CASE 
				WHEN EXISTS (
					SELECT TRUE FROM chat_user WHERE chat_id = @chatID AND user_id = @userID
				) THEN TRUE
				ELSE FALSE
			END as "exists"
	`

	args := pgx.NamedArgs{
		"chatID": chatID,
		"userID": userID,
	}

	rows, err := c.pool.Query(ctx, query, args)
	if err != nil {
		return false, errors.WithStack(err)
	}
	defer rows.Close()

	type BoolStruct struct {
		Exists bool
	}
	answer, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[BoolStruct])
	if err != nil {
		return false, errors.WithStack(err)
	}

	c.logger.Info().Msg(strconv.FormatBool(answer.Exists))

	return answer.Exists, nil
}

//// Update user in db.
//func (u *CRUDUserRepo) Update(ctx context.Context, in *repositories.CRUDUserUpdateIn) error {
//	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
//	params := make([]any, 0)
//
//	query := psql.Update("\"user\"")
//
//	if in.Name != nil {
//		query = query.Set("name", "?")
//		params = append(params, *in.Name)
//	}
//
//	if in.Role != nil {
//		query = query.Set("role", "?")
//		params = append(params, *in.Role)
//	}
//
//	if in.Email != nil {
//		query = query.Set("email", "?")
//		params = append(params, *in.Email)
//	}
//
//	query = query.Set("updated_at", "?")
//	params = append(params, time.Now())
//
//	query = query.Where(squirrel.Eq{"id": "?"})
//	params = append(params, in.ID)
//
//	query = query.Suffix("RETURNING \"id\"")
//
//	sql, _, err := query.ToSql()
//	if err != nil {
//		return errors.WithStack(err)
//	}
//
//	row := u.pool.QueryRow(ctx, sql, params...)
//	var userID int64
//	if err = row.Scan(&userID); errors.Is(err, pgx.ErrNoRows) {
//		return repositories.ErrUserNotFound
//	}
//
//	return err
//}
//
//// Get user from db.
//func (u *CRUDUserRepo) Get(ctx context.Context, in *repositories.CRUDUserGetIn) (*repositories.CRUDUserGetOut, error) {
//	query := `
//	SELECT
//	    id, name, email, role, created_at createdAt, updated_at updatedAt
//	FROM
//	    "user"
//	WHERE
//	    id = @id
//    `
//
//	args := pgx.NamedArgs{
//		"id": in.ID,
//	}
//
//	rows, err := u.pool.Query(ctx, query, args)
//	if err != nil {
//		return &repositories.CRUDUserGetOut{}, errors.WithStack(err)
//	}
//	defer rows.Close()
//
//	out, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[repositories.CRUDUserGetOut])
//	if err != nil {
//		if errors.Is(err, pgx.ErrNoRows) {
//			return nil, repositories.ErrUserNotFound
//		}
//		return nil, errors.WithStack(err)
//	}
//
//	return &out, nil
//}
//
//// Delete user from db.
//func (u *CRUDUserRepo) Delete(ctx context.Context, in *repositories.CRUDUserDeleteIn) error {
//	query := `
//    	DELETE FROM "user"
//		WHERE id = @id
//    `
//
//	args := pgx.NamedArgs{
//		"id": in.ID,
//	}
//
//	_, err := u.pool.Exec(ctx, query, args)
//	if err != nil {
//		return errors.WithStack(err)
//	}
//
//	return nil
//}
