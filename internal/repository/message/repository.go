package message

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mistandok/chat-server/internal/repository"
	"github.com/rs/zerolog"
)

const (
	messageTable = "message"

	chatIDColumn     = "chat_id"
	fromUserIDColumn = "from_user_id"
	textColumn       = "text"
	sentAtColumn     = "sent_at"
	idColumn         = "id"

	messageChatIDFKConstraint     = "fk_chat_id"
	messageFromUserIDFKConstraint = "fk_from_user_id"
)

var _ repository.MessageRepository = (*Repo)(nil)

// Repo ..
type Repo struct {
	pool   *pgxpool.Pool
	logger *zerolog.Logger
}

// NewRepo  get new repo instance.
func NewRepo(pool *pgxpool.Pool, logger *zerolog.Logger) *Repo {
	return &Repo{
		pool:   pool,
		logger: logger,
	}
}
