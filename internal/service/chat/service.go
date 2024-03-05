package chat

import (
	"github.com/mistandok/chat-server/internal/client/db"
	"github.com/mistandok/chat-server/internal/repository"
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
