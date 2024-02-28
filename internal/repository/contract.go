package repository

import (
	"context"

	serviceModel "github.com/mistandok/chat-server/internal/model"
)

// ChatRepository interface for control chat
type ChatRepository interface {
	Create(ctx context.Context, userIDs []serviceModel.UserID) (serviceModel.ChatID, error)
	Delete(ctx context.Context, chatID serviceModel.ChatID) error
	SendMessage(context.Context, serviceModel.Message) error
}
