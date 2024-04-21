package stream

import (
	"sync"

	"github.com/mistandok/chat-server/internal/model"
)

// Chats ..
type Chats struct {
	chats   map[model.ChatID]*Chat
	rwMutex sync.RWMutex
}

// NewChats struct
func NewChats() *Chats {
	return &Chats{
		chats: make(map[model.ChatID]*Chat),
	}
}

// GetOrCreateChat ..
func (c *Chats) GetOrCreateChat(chatID model.ChatID) *Chat {
	var chat *Chat
	var ok bool

	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()

	if chat, ok = c.chats[chatID]; !ok {
		chat = NewChat()
		c.chats[chatID] = chat
	}

	return chat
}
