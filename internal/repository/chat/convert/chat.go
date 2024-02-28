package convert

import (
	serviceModel "github.com/mistandok/chat-server/internal/model"
	"github.com/mistandok/chat-server/internal/repository/chat/model"
)

func ToSliceIntFromSliceServiceUserID(userIDs []serviceModel.UserID) []int64 {
	result := make([]int64, 0, len(userIDs))
	for _, userID := range userIDs {
		result = append(result, int64(userID))
	}

	return result
}

func ToMessageFromServiceMessage(message *serviceModel.Message) *model.Message {
	return &model.Message{
		FromUserID: int64(message.FromUserID),
		Text:       message.Text,
		ToChatID:   int64(message.ToChatID),
		SendTime:   message.SendTime,
	}
}
