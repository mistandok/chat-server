package chat

import (
	"context"
	"fmt"

	"github.com/mistandok/chat-server/internal/service"

	"github.com/mistandok/chat-server/internal/model"
)

// SendMessage ..
func (s *Service) SendMessage(ctx context.Context, message model.Message) error {
	s.logger.Debug().Msg(fmt.Sprintf("попытка послать сообщение: %+v", message))

	userInChat, err := s.chatRepo.IsUserInChat(ctx, message.ToChatID, message.FromUserID)
	if err != nil {
		s.logger.Err(err).Msg(fmt.Sprintf("не удалось сохранить сообщение: %+v", message))
		return fmt.Errorf("ошибка при попытке отправить сообщение: %w", err)
	}

	if !userInChat {
		return service.ErrMsgUserNotInTheChat
	}

	if err := s.messageRepo.Create(ctx, message); err != nil {
		s.logger.Err(err).Msg(fmt.Sprintf("не удалось сохранить сообщение: %+v", message))
		return fmt.Errorf("ошибка при попытке отправить сообщение: %w", err)
	}

	return nil
}
