package chat

import (
	"context"
	"fmt"
)

// Delete ..
func (s *Service) Delete(ctx context.Context, chatID int64) error {
	s.logger.Debug().Msg(fmt.Sprintf("попытка удалить чат: %d", chatID))

	if err := s.chatRepo.Delete(ctx, chatID); err != nil {
		s.logger.Err(err).Msg(fmt.Sprintf("не удалось удалить чат: %d", chatID))
		return fmt.Errorf("ошибка при попытке удаления чата: %w", err)
	}

	return nil
}
