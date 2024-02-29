package service

import (
	"context"

	"github.com/mistandok/chat-server/internal/model"
)

// ChatService ..
type ChatService interface {
	Create(context.Context, []model.UserID) (model.ChatID, error)
	Delete(context.Context, model.ChatID) error
	SendMessage(context.Context, model.Message) error
}
