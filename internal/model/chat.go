package model

import "time"

// Message ..
type Message struct {
	FromUserID int64
	Text       string
	ToChatID   int64
	SendTime   time.Time
}
