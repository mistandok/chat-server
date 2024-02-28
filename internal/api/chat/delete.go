package chat

import (
	"context"

	"github.com/mistandok/chat-server/internal/model"
	"github.com/mistandok/chat-server/pkg/chat_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) Delete(ctx context.Context, request *chat_v1.DeleteRequest) (*emptypb.Empty, error) {
	err := i.chatService.Delete(ctx, model.ChatID(request.Id))
	if err != nil {
		return nil, status.Error(codes.Internal, "прошу понять и простить :(")
	}

	return &emptypb.Empty{}, nil
}
