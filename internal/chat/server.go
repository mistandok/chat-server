package chat

import (
	"context"
	"fmt"

	"github.com/mistandok/chat-server/pkg/chat_v1"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Server chat Server.
type Server struct {
	chat_v1.UnimplementedChatV1Server
	logger *zerolog.Logger
}

// NewServer generate instance for chat Server.
func NewServer(logger *zerolog.Logger) *Server {
	return &Server{
		logger: logger,
	}
}

// Create chat by param.
func (s *Server) Create(_ context.Context, request *chat_v1.CreateRequest) (*chat_v1.CreateResponse, error) {
	s.logger.Debug().Msg(fmt.Sprintf("try create chat: %+v", request))

	return &chat_v1.CreateResponse{Id: 1}, nil
}

// Delete chat by params.
func (s *Server) Delete(_ context.Context, request *chat_v1.DeleteRequest) (*emptypb.Empty, error) {
	s.logger.Debug().Msg(fmt.Sprintf("try delete chat: %+v", request))

	return &emptypb.Empty{}, nil
}

// SendMessage to chat
func (s *Server) SendMessage(_ context.Context, request *chat_v1.SendMessageRequest) (*emptypb.Empty, error) {
	s.logger.Debug().Msg(fmt.Sprintf("try send message to chat: %+v", request))

	return &emptypb.Empty{}, nil
}
