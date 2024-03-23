package repository

import (
	"context"

	serviceModel "github.com/mistandok/chat-server/internal/model"
)

//go:generate ../../bin/mockery --output ./mocks  --inpackage-suffix --all --case snake

// ChatRepository interface for control chat
type ChatRepository interface {
	Create(ctx context.Context) (int64, error)
	Delete(ctx context.Context, chatID int64) error
	IsUserInChat(context.Context, int64, int64) (bool, error)
	LinkChatAndUsers(context.Context, int64, []int64) error
}

// UserRepository interface for control user
type UserRepository interface {
	CreateMass(ctx context.Context, userIDs []int64) error
}

// MessageRepository interface for control message
type MessageRepository interface {
	Create(context.Context, serviceModel.Message) error
}
