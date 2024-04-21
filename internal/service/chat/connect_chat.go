package chat

import (
	"github.com/mistandok/chat-server/internal/model"
	"github.com/mistandok/chat-server/internal/service"
	"github.com/mistandok/chat-server/internal/service/chat/stream"
)

// ConnectChat ..
func (s *Service) ConnectChat(in *model.ConnectChatIn, stream service.StreamChatMessageSender) error {
	userInChat, err := s.chatRepo.IsUserInChat(stream.Context(), in.ChatID, in.UserID)
	if err != nil {
		return err
	}

	if !userInChat {
		return service.ErrUserNotInTheChat
	}

	chatID := model.ChatID(in.ChatID)
	user := &model.User{
		ID:   in.UserID,
		Name: in.UserName,
	}

	chatMessageChannel, ok := s.chatsMessageChannel.GetChannelForChat(chatID)
	if !ok {
		return service.ErrChatNotFound
	}

	chat := s.chats.GetOrCreateChat(chatID)
	chat.SetStreamForUser(user, stream)

	return s.processingChatChannel(chatMessageChannel, chat, user, stream)
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

			for currentUser, userStream := range chat.GetStreamForUsers() {
				if currentUser.ID == user.ID {
					continue
				}
				if err := userStream.Send(message); err != nil {
					return err
				}
			}

		case <-stream.Context().Done():
			chat.DeleteStreamForUser(user)
			return nil
		}
	}
}
