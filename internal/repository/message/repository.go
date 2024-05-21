package message

import (
	"github.com/mistandok/chat-server/internal/repository"
	"github.com/mistandok/platform_common/pkg/db"
	"github.com/rs/zerolog"
)

const (
	messageTable = "message"

	chatIDColumn       = "chat_id"
	fromUserIDColumn   = "from_user_id"
	fromUserNameColumn = "from_user_name"
	textColumn         = "text"
	sentAtColumn       = "sent_at"
	idColumn           = "id"

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
