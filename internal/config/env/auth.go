package env

import (
	"errors"
	"os"

	"github.com/mistandok/chat-server/internal/config"
)

const (
	authHostEnvName = "AUTH_HOST"
	authPortEnvName = "AUTH_PORT"
)

// AuthCfgSearcher searcher for grpc config.
type AuthCfgSearcher struct{}

// NewAuthCfgSearcher get instance for grpc config searcher.
func NewAuthCfgSearcher() *AuthCfgSearcher {
	return &AuthCfgSearcher{}
}

// Get searcher for grpc config.
func (s *AuthCfgSearcher) Get() (*config.AuthConfig, error) {
	host := os.Getenv(authHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("auth host not found")
	}

	port := os.Getenv(authPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("auth port not found")
	}

	return &config.AuthConfig{
		Host: host,
		Port: port,
	}, nil
}
