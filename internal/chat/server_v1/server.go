package server_v1

import (
	"context"
	"fmt"
	"github.com/mistandok/chat-server/internal/repositories"
	"github.com/mistandok/chat-server/pkg/chat_v1"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ChatRepo interface {
	Create(context.Context, *repositories.ChatCreateIn) (*repositories.ChatCreateOut, error)
	Delete(context.Context, *repositories.ChatDeleteIn) error
	SendMessage(context.Context, *repositories.SendMessageIn) error
}

// Server chat Server.
type Server struct {
	chat_v1.UnimplementedChatV1Server
	logger   *zerolog.Logger
	chatRepo ChatRepo
}

// NewServer generate instance for chat Server.
func NewServer(logger *zerolog.Logger, chatRepo ChatRepo) *Server {
	return &Server{
		logger:   logger,
		chatRepo: chatRepo,
	}
}

// Create chat by param.
func (s *Server) Create(ctx context.Context, request *chat_v1.CreateRequest) (*chat_v1.CreateResponse, error) {
	s.logger.Debug().Msg(fmt.Sprintf("try create chat: %+v", request))

	out, err := s.chatRepo.Create(ctx, &repositories.ChatCreateIn{UserIDs: request.UserIDs})
	if err != nil {
		s.logger.Err(err).Msg("не удалось создать чат")
		return &chat_v1.CreateResponse{}, status.Error(codes.Internal, "прошу понять и простить :(")
	}

	return &chat_v1.CreateResponse{Id: out.ID}, nil
}

// Delete chat by params.
func (s *Server) Delete(ctx context.Context, request *chat_v1.DeleteRequest) (*emptypb.Empty, error) {
	s.logger.Debug().Msg(fmt.Sprintf("try delete chat: %+v", request))

	err := s.chatRepo.Delete(ctx, &repositories.ChatDeleteIn{ID: request.Id})
	if err != nil {
		s.logger.Err(err).Msg("не удалось удалить чат")
		return &emptypb.Empty{}, status.Error(codes.Internal, "прошу понять и простить :(")
	}

	return &emptypb.Empty{}, nil
}

// SendMessage to chat
func (s *Server) SendMessage(ctx context.Context, request *chat_v1.SendMessageRequest) (*emptypb.Empty, error) {
	s.logger.Debug().Msg(fmt.Sprintf("try send message to chat: %+v", request))

	err := s.chatRepo.SendMessage(ctx, &repositories.SendMessageIn{
		FromUserID: request.From,
		Message:    request.Text,
		ToChatId:   request.ToChatId,
		SendTime:   request.Timestamp.AsTime(),
	})
	if err != nil {
		switch {
		case errors.Is(err, repositories.ErrChatNotFound) || errors.Is(err, repositories.ErrUserNotFound):
			s.logger.Warn().Msg("не удалось отправить сообщение")
			return &emptypb.Empty{}, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, repositories.ErrUserNotInTheChat):
			s.logger.Warn().Msg("не удалось отправить сообщение")
			return &emptypb.Empty{}, status.Error(codes.InvalidArgument, err.Error())
		default:
			s.logger.Err(err).Msg("не удалось отправить сообщение")
			return &emptypb.Empty{}, status.Error(codes.Internal, "прошу понять и простить :(")
		}
	}

	return &emptypb.Empty{}, nil
}
