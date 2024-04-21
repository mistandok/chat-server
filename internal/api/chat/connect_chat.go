package chat

import (
	"github.com/mistandok/chat-server/internal/convert"
	"github.com/mistandok/chat-server/pkg/chat_v1"
)

// ConnectChat ..
func (i *Implementation) ConnectChat(request *chat_v1.ConnectChatRequest, stream chat_v1.ChatV1_ConnectChatServer) error {
	err := i.chatService.ConnectChat(convert.ToConnectChatInFromDesc(request), convert.NewStreamMessageSender(stream))
	if err != nil {
		return errInternal
	}

	return nil
}
