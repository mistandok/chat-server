package user

import (
	"github.com/mistandok/chat-server/internal/client/db"
	"github.com/mistandok/chat-server/internal/repository"
	"github.com/rs/zerolog"
)

const (
	userTable = "user"

	idColumn = "id"
)

var _ repository.UserRepository = (*Repo)(nil)

// Repo ..
type Repo struct {
	logger *zerolog.Logger
	db     db.Client
}

// NewRepo  get new repo instance.
func NewRepo(logger *zerolog.Logger, dbClient db.Client) *Repo {
	return &Repo{
		logger: logger,
		db:     dbClient,
	}
}
