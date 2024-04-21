package convert

import (
	"github.com/mistandok/chat-server/internal/model"
	desc "github.com/mistandok/chat-server/pkg/chat_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToMessageFromDesc ..
func ToMessageFromDesc(request *desc.SendMessageRequest) *model.Message {
	return &model.Message{
		FromUserID:   request.Message.FromUserId,
		FromUserName: request.Message.FromUserName,
		Text:         request.Message.Text,
		ToChatID:     request.ToChatId,
		SendTime:     request.Message.CreatedAt.AsTime(),
	}
}

// ToConnectChatInFromDesc ..
func ToConnectChatInFromDesc(request *desc.ConnectChatRequest) *model.ConnectChatIn {
	return &model.ConnectChatIn{
		ChatID:   request.ChatId,
		UserID:   request.UserId,
		UserName: request.UserName,
	}
}

// ToDescMessageFromMessage ..
func ToDescMessageFromMessage(message *model.Message) *desc.Message {
	return &desc.Message{
		FromUserId:   message.FromUserID,
		FromUserName: message.FromUserName,
		Text:         message.Text,
		CreatedAt:    timestamppb.New(message.SendTime),
	}
}
