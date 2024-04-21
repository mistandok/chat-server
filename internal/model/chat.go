package model

import "time"

// Message ..
type Message struct {
	FromUserID   int64
	FromUserName string
	Text         string
	ToChatID     int64
	SendTime     time.Time
}
