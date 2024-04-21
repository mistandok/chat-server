package stream

import (
	"sync"

	"github.com/mistandok/chat-server/internal/model"
	"github.com/mistandok/chat-server/internal/service"
)

// Chat ..
type Chat struct {
	streamForUsers map[model.User]service.StreamChatMessageSender
	rwMutex        sync.RWMutex
}

// NewChat ..
func NewChat() *Chat {
	return &Chat{
		streamForUsers: make(map[model.User]service.StreamChatMessageSender),
	}
}

// SetStreamForUser ..
func (c *Chat) SetStreamForUser(user *model.User, stream service.StreamChatMessageSender) {
	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()

	c.streamForUsers[*user] = stream
}

// GetStreamForUsers ..
func (c *Chat) GetStreamForUsers() map[model.User]service.StreamChatMessageSender {
	c.rwMutex.RLock()
	defer c.rwMutex.RUnlock()

	streamForUsers := make(map[model.User]service.StreamChatMessageSender, len(c.streamForUsers))
	for user, stream := range c.streamForUsers {
		streamForUsers[user] = stream
	}

	return streamForUsers
}

// GetStreamForUser ..
func (c *Chat) GetStreamForUser(user *model.User) (service.StreamChatMessageSender, bool) {
	c.rwMutex.RLock()
	defer c.rwMutex.RUnlock()

	var stream service.StreamChatMessageSender
	var ok bool

	stream, ok = c.streamForUsers[*user]

	return stream, ok
}

// DeleteStreamForUser ..
func (c *Chat) DeleteStreamForUser(user *model.User) {
	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()

	delete(c.streamForUsers, *user)
}
