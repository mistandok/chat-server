package convert

import (
	serviceModel "github.com/mistandok/chat-server/internal/model"
)

// ToSliceIntFromSliceServiceUserID ..
func ToSliceIntFromSliceServiceUserID(userIDs []serviceModel.UserID) []int64 {
	result := make([]int64, 0, len(userIDs))
	for _, userID := range userIDs {
		result = append(result, int64(userID))
	}

	return result
}
