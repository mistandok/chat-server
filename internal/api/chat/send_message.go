package chat

import (
	"context"

	"github.com/mistandok/chat-server/internal/convert"
	"github.com/mistandok/chat-server/internal/repository"
	"github.com/mistandok/chat-server/pkg/chat_v1"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) SendMessage(ctx context.Context, request *chat_v1.SendMessageRequest) (*emptypb.Empty, error) {
	err := i.chatService.SendMessage(ctx, *convert.ToMessageFromDesc(request))
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrChatNotFound) || errors.Is(err, repository.ErrUserNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, repository.ErrUserNotInTheChat):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, "прошу понять и простить :(")
		}
	}

	return &emptypb.Empty{}, nil
}
