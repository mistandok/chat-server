package service

import (
	"context"

	"github.com/mistandok/chat-server/internal/model"
)

// ChatService ..
type ChatService interface {
	Create(context.Context, []int64) (int64, error)
	Delete(context.Context, int64) error
	SendMessage(context.Context, model.Message) error
}
