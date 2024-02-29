package user

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mistandok/chat-server/internal/repository"
	"github.com/rs/zerolog"
)

const (
	userTable = "user"

	userIDColumn = "user_id"
	idColumn     = "id"
)

var _ repository.UserRepository = (*Repo)(nil)

// Repo ..
type Repo struct {
	pool   *pgxpool.Pool
	logger *zerolog.Logger
}

// NewRepo  get new repo instance.
func NewRepo(pool *pgxpool.Pool, logger *zerolog.Logger) *Repo {
	return &Repo{
		pool:   pool,
		logger: logger,
	}
}
