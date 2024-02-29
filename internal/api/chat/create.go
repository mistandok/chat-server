package chat

import (
	"context"

	"github.com/mistandok/chat-server/internal/convert"
	"github.com/mistandok/chat-server/pkg/chat_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Create ..
func (i *Implementation) Create(ctx context.Context, request *chat_v1.CreateRequest) (*chat_v1.CreateResponse, error) {
	chatID, err := i.chatService.Create(ctx, convert.ToSliceUserIDsFromSliceIntDesc(request.UserIDs))
	if err != nil {
		return nil, status.Error(codes.Internal, "прошу понять и простить :(")
	}

	return &chat_v1.CreateResponse{Id: int64(chatID)}, nil
}
