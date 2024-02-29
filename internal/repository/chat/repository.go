package chat

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mistandok/chat-server/internal/client/db"
	"github.com/mistandok/chat-server/internal/repository"
	"github.com/rs/zerolog"
)

const (
	chatTable     = "chat"
	chatUserTable = "chat_user"

	userIDColumn = "user_id"
	chatIDColumn = "chat_id"
	idColumn     = "id"
)

var _ repository.ChatRepository = (*Repo)(nil)

// Repo ..
type Repo struct {
	pool   *pgxpool.Pool
	logger *zerolog.Logger
	db     db.Client
}

// NewRepo  get new repo instance.
func NewRepo(pool *pgxpool.Pool, logger *zerolog.Logger, dbClient db.Client) *Repo {
	return &Repo{
		pool:   pool,
		logger: logger,
		db:     dbClient,
	}
}
