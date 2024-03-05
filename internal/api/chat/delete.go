package chat

import (
	"context"

	"github.com/mistandok/chat-server/pkg/chat_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Delete ..
func (i *Implementation) Delete(ctx context.Context, request *chat_v1.DeleteRequest) (*emptypb.Empty, error) {
	err := i.chatService.Delete(ctx, request.Id)
	if err != nil {
		return nil, errInternal
	}

	return &emptypb.Empty{}, nil
}
