package access

import (
	"context"
	"fmt"

	"github.com/mistandok/chat-server/internal/client"
	"github.com/mistandok/chat-server/pkg/auth_v1"
	"github.com/rs/zerolog"
)

var _ client.AccessClient = (*ClientFacade)(nil)

// ClientFacade ..
type ClientFacade struct {
	logger *zerolog.Logger
	client auth_v1.AccessV1Client
}

// NewClientFacade ..
func NewClientFacade(logger *zerolog.Logger, client auth_v1.AccessV1Client) *ClientFacade {
	return &ClientFacade{
		logger: logger,
		client: client,
	}
}

// Check ..
func (a *ClientFacade) Check(ctx context.Context, address string) bool {
	a.logger.Debug().Msg(fmt.Sprintf("запрос доступа на address: %s", address))
	_, err := a.client.Check(ctx, &auth_v1.CheckRequest{Address: address})
	if err != nil {
		a.logger.Warn().Err(err).Msg("ошибка при запросе доступа")
	}
	return err == nil
}
