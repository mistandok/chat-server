package model

import "time"

type ChatID int64

type UserID int64

type Message struct {
	FromUserID UserID
	Text       string
	ToChatID   ChatID
	SendTime   time.Time
}
