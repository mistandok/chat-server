package model

import "time"

// ChatID ..
type ChatID int64

// UserID ..
type UserID int64

// Message ..
type Message struct {
	FromUserID UserID
	Text       string
	ToChatID   ChatID
	SendTime   time.Time
}
