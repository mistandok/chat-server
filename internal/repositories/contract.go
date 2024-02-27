package repositories

import (
	"context"
	"time"
)

// ChatCreateIn input data for create chat.
type ChatCreateIn struct {
	UserIDs []int64 `json:"userIDs"`
}

// ChatCreateOut out data after create chat.
type ChatCreateOut struct {
	ID int64 `json:"id"`
}

// ChatDeleteIn input data for delete chat.
type ChatDeleteIn struct {
	ID int64 `json:"id"`
}

// SendMessageIn input data for send message.
type SendMessageIn struct {
	FromUserID int64     `json:"from_user_id"`
	Message    string    `json:"message"`
	ToChatID   int64     `json:"to_chat_id"`
	SendTime   time.Time `json:"send_time"`
}

// ChatRepository interface for control chat
type ChatRepository interface {
	Create(context.Context, *ChatCreateIn) (*ChatCreateOut, error)
	Delete(context.Context, *ChatDeleteIn) error
	SendMessage(context.Context, *SendMessageIn) error
}
