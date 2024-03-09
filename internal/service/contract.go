package service

import (
	"context"

	"github.com/mistandok/chat-server/internal/model"
)

//go:generate ../../bin/mockery --output ./mocks  --inpackage-suffix --all

// ChatService ..
type ChatService interface {
	Create(context.Context, []int64) (int64, error)
	Delete(context.Context, int64) error
	SendMessage(context.Context, model.Message) error
}
