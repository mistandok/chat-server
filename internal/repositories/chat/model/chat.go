package model

import (
	"time"
)

type Message struct {
	FromUserID int64     `db:"from_user_id"`
	Message    string    `db:"message"`
	ToChatID   int64     `db:"to_chat_id"`
	SendTime   time.Time `db:"send_time"`
}
