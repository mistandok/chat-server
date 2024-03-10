package tests

import (
	"time"

	"github.com/mistandok/chat-server/internal/model"
	"github.com/mistandok/chat-server/pkg/chat_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func sendMessageRequestWithSetup(fromUserID int64, toChatID int64, onTime time.Time) *chat_v1.SendMessageRequest {
	return &chat_v1.SendMessageRequest{
		From:      fromUserID,
		Text:      "msg",
		Timestamp: timestamppb.New(onTime),
		ToChatId:  toChatID,
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
