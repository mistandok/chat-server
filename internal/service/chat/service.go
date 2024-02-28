package chat

import (
	"github.com/mistandok/chat-server/internal/repository"
	"github.com/rs/zerolog"
)

type Service struct {
	chatRepo repository.ChatRepository
	logger   *zerolog.Logger
}

func NewService(chatRepo repository.ChatRepository, logger *zerolog.Logger) *Service {
	return &Service{
		chatRepo: chatRepo,
		logger:   logger,
	}
}
