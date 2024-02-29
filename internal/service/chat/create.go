package chat

import (
	"context"
	"fmt"

	"github.com/mistandok/chat-server/internal/model"
)

// Create ..
func (s *Service) Create(ctx context.Context, userIDs []model.UserID) (model.ChatID, error) {
	s.logger.Debug().Msg("попытка создать чат")
	chatID, err := s.chatRepo.Create(ctx)
	if err != nil {
		s.logger.Err(err).Msg("не удалось создать чат")
		return 0, fmt.Errorf("ошибка при попытке создания чата: %w", err)
	}

	err = s.userRepo.CreateMass(ctx, userIDs)
	if err != nil {
		s.logger.Err(err).Msg("не удалось создать пользователей")
		return 0, fmt.Errorf("ошибка при попытке создания чата: %w", err)
	}

	err = s.chatRepo.LinkChatAndUsers(ctx, chatID, userIDs)
	if err != nil {
		s.logger.Err(err).Msg("не удалось связать пользователей с чатом")
		return 0, fmt.Errorf("ошибка при попытке создания чата: %w", err)
	}

	return chatID, nil
}
