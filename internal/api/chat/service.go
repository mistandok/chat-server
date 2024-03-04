package chat

import (
	"errors"

	"github.com/mistandok/chat-server/internal/service"
	"github.com/mistandok/chat-server/pkg/chat_v1"
)

const msgInternalError = "что-то пошло не так, мы уже работаем над решением проблемы"

var errInternal = errors.New(msgInternalError)

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
