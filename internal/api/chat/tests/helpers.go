package tests

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"

	"github.com/mistandok/chat-server/internal/model"
	"github.com/mistandok/chat-server/pkg/chat_v1"
)

func sendMessageRequestWithSetup(fromUserID int64, toChatID int64, onTime time.Time) *chat_v1.SendMessageRequest {
	return &chat_v1.SendMessageRequest{
		Message: &chat_v1.Message{
			FromUserID: fromUserID,
			Text:       "msg",
			CreatedAt:  timestamppb.New(onTime),
		},
		ToChatId: toChatID,
	}
}

func messageWithSetup(fromUserID int64, toChatID int64, onTime time.Time) *model.Message {
	return &model.Message{
		FromUserID: fromUserID,
		Text:       "msg",
		ToChatID:   toChatID,
		SendTime:   onTime,
	}
}
