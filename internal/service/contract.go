package service

import (
	"context"

	"github.com/mistandok/chat-server/internal/model"
)

//go:generate ../../bin/mockery --output ./mocks  --inpackage-suffix --all --case snake

// ChatService ..
type ChatService interface {
	Create(context.Context, []int64) (int64, error)
	Delete(context.Context, int64) error
	SendMessage(context.Context, model.Message) error
	ConnectChat(*model.ConnectChatIn, StreamChatMessageSender) error
}

// StreamChatMessageSender interface for send message in chat stream
type StreamChatMessageSender interface {
	Send(*model.Message) error
	Context() context.Context
}
