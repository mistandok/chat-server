package client

import (
	"context"
)

type AccessClient interface {
	Check(ctx context.Context, address string) bool
}
