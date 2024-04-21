package convert

import (
	"context"

	"github.com/mistandok/chat-server/internal/model"
	"github.com/mistandok/chat-server/internal/service"
	"github.com/mistandok/chat-server/pkg/chat_v1"
)

var _ service.StreamChatMessageSender = (*StreamMessageSender)(nil)

// StreamMessageSender ..
type StreamMessageSender struct {
	chatServer chat_v1.ChatV1_ConnectChatServer
}

// Send ..
func (s *StreamMessageSender) Send(message *model.Message) error {
	return s.chatServer.Send(ToDescMessageFromMessage(message))
}

// Context ..
func (s *StreamMessageSender) Context() context.Context {
	return s.chatServer.Context()
}

// NewStreamMessageSender ..
func NewStreamMessageSender(chatServer chat_v1.ChatV1_ConnectChatServer) service.StreamChatMessageSender {
	return &StreamMessageSender{chatServer: chatServer}
}
