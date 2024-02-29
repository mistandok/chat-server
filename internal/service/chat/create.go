package chat

import (
	"context"
	"fmt"

	"github.com/mistandok/chat-server/internal/model"
)

// Create ..
func (s *Service) Create(ctx context.Context, userIDs []model.UserID) (model.ChatID, error) {
	s.logger.Debug().Msg("попытка создать чат")

	var chatID model.ChatID
	var err error

	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var txErr error
		chatID, txErr = s.chatRepo.Create(ctx)
		if txErr != nil {
			s.logger.Err(txErr).Msg("не удалось создать чат")
			return fmt.Errorf("ошибка при попытке создания чата: %w", txErr)
		}

		txErr = s.userRepo.CreateMass(ctx, userIDs)
		if err != nil {
			s.logger.Err(txErr).Msg("не удалось создать пользователей")
			return fmt.Errorf("ошибка при попытке создания пользователей: %w", txErr)
		}

		txErr = s.chatRepo.LinkChatAndUsers(ctx, chatID, userIDs)
		if txErr != nil {
			s.logger.Err(txErr).Msg("не удалось связать пользователей с чатом")
			return fmt.Errorf("ошибка при попытке свящать пользователей с чатом: %w", txErr)
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return chatID, nil
}
