package convert

import (
	serviceModel "github.com/mistandok/chat-server/internal/model"
	"github.com/mistandok/chat-server/internal/repository/message/model"
)

// ToMessageFromServiceMessage ..
func ToMessageFromServiceMessage(message *serviceModel.Message) *model.Message {
	return &model.Message{
		FromUserID: int64(message.FromUserID),
		Text:       message.Text,
		ToChatID:   int64(message.ToChatID),
		SendTime:   message.SendTime,
	}
}
