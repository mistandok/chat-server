package model

import "time"

// Message ..
type Message struct {
	FromUserID   int64     `db:"from_user_id"`
	FromUserName string    `db:"from_user_name"`
	Text         string    `db:"text"`
	ToChatID     int64     `db:"to_chat_id"`
	SendTime     time.Time `db:"send_time"`
}
