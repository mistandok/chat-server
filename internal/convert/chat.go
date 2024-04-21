package convert

import (
	"github.com/mistandok/chat-server/internal/model"
	desc "github.com/mistandok/chat-server/pkg/chat_v1"
)

// ToMessageFromDesc ..
func ToMessageFromDesc(request *desc.SendMessageRequest) *model.Message {
	return &model.Message{
		FromUserID:   request.Message.FromUserID,
		FromUserName: request.Message.FromUserName,
		Text:         request.Message.Text,
		ToChatID:     request.ToChatId,
		SendTime:     request.Message.CreatedAt.AsTime(),
	}
}
