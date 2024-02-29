package repository

import (
	"context"

	serviceModel "github.com/mistandok/chat-server/internal/model"
)

// ChatRepository interface for control chat
type ChatRepository interface {
	Create(ctx context.Context) (serviceModel.ChatID, error)
	Delete(ctx context.Context, chatID serviceModel.ChatID) error
	IsUserInChat(context.Context, serviceModel.ChatID, serviceModel.UserID) (bool, error)
	LinkChatAndUsers(context.Context, serviceModel.ChatID, []serviceModel.UserID) error
}

// UserRepository interface for control user
type UserRepository interface {
	CreateMass(ctx context.Context, userIDs []serviceModel.UserID) error
}

// MessageRepository interface for control message
type MessageRepository interface {
	Create(context.Context, serviceModel.Message) error
}
