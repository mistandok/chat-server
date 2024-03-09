package chat

import (
	"github.com/mistandok/chat-server/internal/repository"
	"github.com/mistandok/platform_common/pkg/db"
	"github.com/rs/zerolog"
)

// Service ..
type Service struct {
	logger      *zerolog.Logger
	txManager   db.TxManager
	chatRepo    repository.ChatRepository
	userRepo    repository.UserRepository
	messageRepo repository.MessageRepository
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
		chatRepo:    chatRepo,
		userRepo:    userRepo,
		messageRepo: messageRepo,
		logger:      logger,
		txManager:   txManager,
	}
}
