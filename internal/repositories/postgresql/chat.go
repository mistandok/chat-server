package postgresql

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mistandok/chat-server/internal/repositories"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const (
	chatTable     = "chat"
	chatUserTable = "chat_user"
	userTable     = "user"
	messageTable  = "message"

	userIdColumn     = "user_id"
	chatIdColumn     = "chat_id"
	fromUserIdColumn = "from_user_id"
	messageColumn    = "message"
	sentAtColumn     = "sent_at"
	idColumn         = "id"

	messageChatIdFKConstraint     = "fk_chat_id"
	messageFromUserIdFKConstraint = "fk_from_user_id"
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
		return nil, err
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			err = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
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
	queryFormat := `DELETE FROM %s WHERE %s = @%s`
	query := fmt.Sprintf(queryFormat, chatTable, idColumn, idColumn)

	args := pgx.NamedArgs{
		idColumn: in.ID,
	}

	_, err := c.pool.Exec(ctx, query, args)

	return err
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

	queryFormat := `
		INSERT INTO %s (%s, %s, %s, %s)
		VALUES (@%s, @%s, @%s, @%s)
    `
	query := fmt.Sprintf(
		queryFormat,
		messageTable,
		fromUserIdColumn, chatIdColumn, messageColumn, sentAtColumn,
		fromUserIdColumn, chatIdColumn, messageColumn, sentAtColumn,
	)

	args := pgx.NamedArgs{
		fromUserIdColumn: in.FromUserID,
		chatIdColumn:     in.ToChatID,
		messageColumn:    in.Message,
		sentAtColumn:     in.SendTime,
	}

	_, err = c.pool.Exec(ctx, query, args)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch {
			case pgErr.ConstraintName == messageChatIdFKConstraint:
				return repositories.ErrChatNotFound
			case pgErr.ConstraintName == messageFromUserIdFKConstraint:
				return repositories.ErrUserNotFound
			}
		}
		return err
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

func (c *ChatRepo) createUsers(ctx context.Context, tx pgx.Tx, userIDs []int64) error {
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
		pgx.Identifier{chatUserTable},
		[]string{chatIdColumn, userIdColumn},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *ChatRepo) isUserInChat(ctx context.Context, chatID int64, userID int64) (bool, error) {
	queryFormat := `
		SELECT TRUE
		FROM %s
		WHERE %s = @%s AND %s = @%s
    `
	query := fmt.Sprintf(
		queryFormat,
		chatUserTable,
		chatIdColumn, chatIdColumn, userIdColumn, userIdColumn,
	)

	c.logger.Info().Msg(query)

	args := pgx.NamedArgs{
		chatIdColumn: chatID,
		userIdColumn: userID,
	}

	rows, err := c.pool.Query(ctx, query, args)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	_, err = pgx.CollectOneRow(rows, pgx.RowTo[bool])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
