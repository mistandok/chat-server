package chat

import (
	"github.com/mistandok/chat-server/internal/repository"
	"github.com/rs/zerolog"
)

// Service ..
type Service struct {
	chatRepo    repository.ChatRepository
	userRepo    repository.UserRepository
	messageRepo repository.MessageRepository
	logger      *zerolog.Logger
}

// NewService ..
func NewService(
	chatRepo repository.ChatRepository,
	userRepo repository.UserRepository,
	messageRepo repository.MessageRepository,
	logger *zerolog.Logger,
) *Service {
	return &Service{
		chatRepo:    chatRepo,
		userRepo:    userRepo,
		messageRepo: messageRepo,
		logger:      logger,
	}
}
