package chat

import (
	"context"
	"errors"

	"github.com/mistandok/chat-server/internal/service"

	"github.com/mistandok/chat-server/internal/convert"
	"github.com/mistandok/chat-server/internal/repository"
	"github.com/mistandok/chat-server/pkg/chat_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// SendMessage ..
func (i *Implementation) SendMessage(ctx context.Context, request *chat_v1.SendMessageRequest) (*emptypb.Empty, error) {
	err := i.chatService.SendMessage(ctx, *convert.ToMessageFromDesc(request))
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrChatNotFound) || errors.Is(err, repository.ErrUserNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, service.ErrMsgUserNotInTheChat):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, errInternal
		}
	}

	return &emptypb.Empty{}, nil
}
