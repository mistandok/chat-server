package utils

import (
	"context"
	"errors"
	"strings"

	"google.golang.org/grpc/metadata"
)

// RotateBearerAuthFromIncomingToOutgoingCtx переменещие auth данных из входящего в исходящий контекст
func RotateBearerAuthFromIncomingToOutgoingCtx(ctx context.Context) (context.Context, error) {
	incomingMd, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("метаданные не переданы")
	}

	authHeader, ok := incomingMd["authorization"]
	if !ok || len(authHeader) == 0 {
		return nil, errors.New("в header не представлена authorization")
	}

	if !strings.HasPrefix(authHeader[0], "Bearer ") {
		return nil, errors.New("некоректный формат authorization в header")
	}

	accessToken := strings.TrimPrefix(authHeader[0], "Bearer ")
	outgoingMd := metadata.New(map[string]string{"Authorization": "Bearer " + accessToken})

	return metadata.NewOutgoingContext(ctx, outgoingMd), nil
}
