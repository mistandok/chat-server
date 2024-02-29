package chat

import (
	"github.com/mistandok/chat-server/internal/service"
	"github.com/mistandok/chat-server/pkg/chat_v1"
)

// Implementation ..
type Implementation struct {
	chat_v1.UnimplementedChatV1Server
	chatService service.ChatService
}

// NewImplementation ..
func NewImplementation(chatService service.ChatService) *Implementation {
	return &Implementation{
		chatService: chatService,
	}
}
