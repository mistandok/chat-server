package client

import (
	"context"
)

// AccessClient ..
type AccessClient interface {
	Check(ctx context.Context, address string) bool
}
