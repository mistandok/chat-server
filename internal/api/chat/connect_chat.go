package chat

import (
	"github.com/mistandok/chat-server/pkg/chat_v1"
)

func (i *Implementation) ConnectChat(request *chat_v1.ConnectChatRequest, stream chat_v1.ChatV1_ConnectChatServer) error {
	return nil
}
