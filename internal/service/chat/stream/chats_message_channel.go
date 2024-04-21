package stream

import (
	"errors"
	"sync"

	"github.com/mistandok/chat-server/internal/model"
)

// ChatsMessageChannel ..
type ChatsMessageChannel struct {
	channels map[model.ChatID]chan *model.Message
	rwMutex  sync.RWMutex
}

// NewChatsMessageChannel ..
func NewChatsMessageChannel() *ChatsMessageChannel {
	return &ChatsMessageChannel{
		channels: make(map[model.ChatID]chan *model.Message),
	}
}

// GetChannelForChat ..
func (c *ChatsMessageChannel) GetChannelForChat(chatID model.ChatID) (chan *model.Message, bool) {
	var chatMessageChannel chan *model.Message
	var ok bool

	c.rwMutex.RLock()
	defer c.rwMutex.RUnlock()

	chatMessageChannel, ok = c.channels[chatID]

	return chatMessageChannel, ok
}

// InitMsgChannelForChat ..
func (c *ChatsMessageChannel) InitMsgChannelForChat(chatID model.ChatID, bufferSize int) {
	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()

	c.channels[chatID] = make(chan *model.Message, bufferSize)
}

// SendMessageToChannelForChat ..
func (c *ChatsMessageChannel) SendMessageToChannelForChat(chatID model.ChatID, message *model.Message) error {
	channel, ok := c.GetChannelForChat(chatID)
	if !ok {
		return errors.New("не найден канал для чата")
	}

	channel <- message

	return nil
}
