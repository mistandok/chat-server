package interceptor

import (
	"context"
	"errors"
	"github.com/mistandok/chat-server/internal/client"
	"github.com/mistandok/chat-server/internal/utils"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

// AccessCheckInterceptor интерцептор для проверки токена на валидности и получении авторизационных прав.
type AccessCheckInterceptor struct {
	logger             *zerolog.Logger
	accessClientFacade client.AccessClient
}

func NewAccessCheckInterceptor(logger *zerolog.Logger, accessClientFacade client.AccessClient) *AccessCheckInterceptor {
	return &AccessCheckInterceptor{
		logger:             logger,
		accessClientFacade: accessClientFacade,
	}
}

// Get реализация интерцептора
func (a *AccessCheckInterceptor) Get(ctx context.Context, req interface{}, serverInfo *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	a.logger.Debug().Msg("AccessCheckInterceptor: попытка достать access token из ctx")
	ctx, err := utils.RotateBearerAuthFromIncomingToOutgoingCtx(ctx)
	if err != nil {
		a.logger.Error().Err(err)
		return nil, err
	}

	if a.accessClientFacade.Check(ctx, serverInfo.FullMethod) {
		return handler(ctx, req)
	}

	return nil, errors.New("доступ запрещен")
}
