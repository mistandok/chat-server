package chat

import (
	"context"
	"fmt"

	"github.com/mistandok/chat-server/internal/model"
)

func (s *Service) Create(ctx context.Context, userIDs []model.UserID) (model.ChatID, error) {
	s.logger.Debug().Msg("попытка создать чат")
	chatID, err := s.chatRepo.Create(ctx, userIDs)
	if err != nil {
		s.logger.Err(err).Msg("не удалось создать чат")
		return 0, fmt.Errorf("ошибка при попытке создания чата: %w", err)
	}

	return chatID, nil
}
