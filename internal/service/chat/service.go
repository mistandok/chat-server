package chat

import (
	"github.com/mistandok/chat-server/internal/repository"
	"github.com/mistandok/chat-server/internal/service"
	"github.com/mistandok/chat-server/internal/service/chat/stream"
	"github.com/mistandok/platform_common/pkg/db"
	"github.com/rs/zerolog"
)

var _ service.ChatService = (*Service)(nil)

// Service ..
type Service struct {
	logger      *zerolog.Logger
	txManager   db.TxManager
	chatRepo    repository.ChatRepository
	userRepo    repository.UserRepository
	messageRepo repository.MessageRepository

	chats               *stream.Chats
	chatsMessageChannel *stream.ChatsMessageChannel
}

// NewService ..
func NewService(
	logger *zerolog.Logger,
	txManager db.TxManager,
	chatRepo repository.ChatRepository,
	userRepo repository.UserRepository,
	messageRepo repository.MessageRepository,
) *Service {
	return &Service{
		chatRepo:            chatRepo,
		userRepo:            userRepo,
		messageRepo:         messageRepo,
		logger:              logger,
		txManager:           txManager,
		chats:               stream.NewChats(),
		chatsMessageChannel: stream.NewChatsMessageChannel(),
	}
}
