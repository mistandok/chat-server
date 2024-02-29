package message

import (
	"github.com/mistandok/chat-server/internal/client/db"
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
	logger *zerolog.Logger
	db     db.Client
}

// NewRepo  get new repo instance.
func NewRepo(logger *zerolog.Logger, dbClient db.Client) *Repo {
	return &Repo{
		logger: logger,
		db:     dbClient,
	}
}
