package chat

import (
	"context"
	"fmt"

	"github.com/mistandok/chat-server/internal/model"
)

func (s *Service) SendMessage(ctx context.Context, message model.Message) error {
	s.logger.Debug().Msg(fmt.Sprintf("попытка послать сообщение: %+v", message))
	if err := s.chatRepo.SendMessage(ctx, message); err != nil {
		s.logger.Err(err).Msg(fmt.Sprintf("не удалось послать сообщение сообщение: %+v", message))
		return fmt.Errorf("ошибка при попытке отправить сообщение: %w", err)
	}

	return nil
}
