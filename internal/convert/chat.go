package convert

import (
	"github.com/mistandok/chat-server/internal/model"
	desc "github.com/mistandok/chat-server/pkg/chat_v1"
)

// ToSliceUserIDsFromSliceIntDesc ..
func ToSliceUserIDsFromSliceIntDesc(userIDs []int64) []model.UserID {
	result := make([]model.UserID, 0, len(userIDs))
	for _, userID := range userIDs {
		result = append(result, model.UserID(userID))
	}

	return result
}

// ToMessageFromDesc ..
func ToMessageFromDesc(request *desc.SendMessageRequest) *model.Message {
	return &model.Message{
		FromUserID: model.UserID(request.From),
		Text:       request.Text,
		ToChatID:   model.ChatID(request.ToChatId),
		SendTime:   request.Timestamp.AsTime(),
	}
}
