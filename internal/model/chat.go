package model

import (
	"time"
)

// Message ..
type Message struct {
	FromUserID   int64
	FromUserName string
	Text         string
	ToChatID     int64
	SendTime     time.Time
}

// ConnectChatIn ..
type ConnectChatIn struct {
	ChatID   int64
	UserID   int64
	UserName string
}

// ChatID ..
type ChatID int64

// User ..
type User struct {
	ID   int64
	Name string
}
