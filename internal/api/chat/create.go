package chat

import (
	"context"

	"github.com/mistandok/chat-server/pkg/chat_v1"
)

// Create ..
func (i *Implementation) Create(ctx context.Context, request *chat_v1.CreateRequest) (*chat_v1.CreateResponse, error) {
	chatID, err := i.chatService.Create(ctx, request.UserIDs)
	if err != nil {
		return nil, errInternal
	}

	return &chat_v1.CreateResponse{Id: chatID}, nil
}
