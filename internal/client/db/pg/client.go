package pg

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mistandok/chat-server/internal/client/db"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type pgClient struct {
	masterDBC db.DB
	logger    *zerolog.Logger
}

// New новый клиент для работы с Postgres
func New(ctx context.Context, dsn string, logger *zerolog.Logger) (db.Client, error) {
	pgxConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, errors.Errorf("ошибка при формировании конфига для pgxpool: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, errors.Errorf("ошибка при подключении к БД: %v", err)
	}

	return &pgClient{
		masterDBC: NewDB(pool, logger),
		logger:    logger,
	}, nil
}

// DB доступ к интерфейсу базы данных
func (c *pgClient) DB() db.DB {
	return c.masterDBC
}

// Close закрытие соединений
func (c *pgClient) Close() error {
	if c.masterDBC != nil {
		c.masterDBC.Close()
	}

	return nil
}
