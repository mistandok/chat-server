package chat

import (
	"context"

	"github.com/mistandok/chat-server/internal/model"
	"github.com/mistandok/chat-server/internal/service"
	"github.com/mistandok/chat-server/internal/service/chat/stream"
)

// ConnectChat ..
func (s *Service) ConnectChat(in *model.ConnectChatIn, stream service.StreamChatMessageSender) error {
	s.logger.Debug().Msg("попытка установить соединение с чатом")

	userInChat, err := s.chatRepo.IsUserInChat(stream.Context(), in.ChatID, in.UserID)
	if err != nil {
		s.logger.Err(err).Msg("не удалось проверить, что пользователь находится в чате")
		return err
	}

	if !userInChat {
		err = s.createUserIfIsNeedAndLinkWithChat(stream.Context(), in.ChatID, in.UserID)
		if err != nil {
			return err
		}
	}

	chatID := model.ChatID(in.ChatID)
	user := &model.User{
		ID:   in.UserID,
		Name: in.UserName,
	}

	chatMessageChannel, ok := s.chatsMessageChannel.GetChannelForChat(chatID)
	if !ok {
		chatMessageChannel = s.chatsMessageChannel.InitMsgChannelForChat(chatID, 100)
	}

	chat := s.chats.GetOrCreateChat(chatID)
	chat.SetStreamForUser(user, stream)

	return s.processingChatChannel(chatMessageChannel, chat, user, stream)
}

func (s *Service) createUserIfIsNeedAndLinkWithChat(ctx context.Context, chatID int64, userID int64) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var txErr error
		txErr = s.userRepo.CreateMass(ctx, []int64{userID})
		if txErr != nil {
			s.logger.Err(txErr).Msg("не удалось создать пользовтеля")
			return txErr
		}

		txErr = s.chatRepo.LinkChatAndUsers(ctx, chatID, []int64{userID})
		if txErr != nil {
			s.logger.Err(txErr).Msg("не удалось связать пользователя с чатом")
		}

		return nil
	})

	return err
}

func (s *Service) processingChatChannel(
	chatChannel chan *model.Message,
	chat *stream.Chat,
	user *model.User,
	stream service.StreamChatMessageSender,
) error {
	for {
		select {
		case message, okChannel := <-chatChannel:
			if !okChannel {
				return nil
			}

			for _, userStream := range chat.GetStreamForUsers() {
				if err := userStream.Send(message); err != nil {
					s.logger.Err(err).Msg("не удалось отправить сообщение в user stream")
					return err
				}
			}

		case <-stream.Context().Done():
			chat.DeleteStreamForUser(user)
			return nil
		}
	}
}
