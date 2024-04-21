package convert

import (
	serviceModel "github.com/mistandok/chat-server/internal/model"
	"github.com/mistandok/chat-server/internal/repository/message/model"
)

// ToMessageFromServiceMessage ..
func ToMessageFromServiceMessage(message *serviceModel.Message) *model.Message {
	return &model.Message{
		FromUserID:   message.FromUserID,
		FromUserName: message.FromUserName,
		Text:         message.Text,
		ToChatID:     message.ToChatID,
		SendTime:     message.SendTime,
	}
}
