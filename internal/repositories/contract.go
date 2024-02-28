package repositories

import (
	"context"
	"github.com/mistandok/chat-server/internal/repositories/chat/model"
)

// ChatRepository interface for control chat
type ChatRepository interface {
	Create(context.Context, []int64) (int64, error)
	Delete(context.Context, int64) error
	SendMessage(context.Context, *model.Message) error
}
