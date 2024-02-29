package app

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mistandok/chat-server/internal/api/chat"
	"github.com/mistandok/chat-server/internal/closer"
	"github.com/mistandok/chat-server/internal/config"
	"github.com/mistandok/chat-server/internal/config/env"
	"github.com/mistandok/chat-server/internal/repository"
	chatRepository "github.com/mistandok/chat-server/internal/repository/chat"
	messageRepository "github.com/mistandok/chat-server/internal/repository/message"
	userRepository "github.com/mistandok/chat-server/internal/repository/user"
	"github.com/mistandok/chat-server/internal/service"
	chatService "github.com/mistandok/chat-server/internal/service/chat"
	"github.com/rs/zerolog"
	"log"
	"os"
)

type serviceProvider struct {
	pgConfig   *config.PgConfig
	grpcConfig *config.GRPCConfig
	logger     *zerolog.Logger

	pool *pgxpool.Pool

	chatRepo    repository.ChatRepository
	userRepo    repository.UserRepository
	messageRepo repository.MessageRepository

	chatService service.ChatService

	chatImpl *chat.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

// PgConfig ..
func (s *serviceProvider) PgConfig() *config.PgConfig {
	if s.pgConfig == nil {
		cfgSearcher := env.NewPgCfgSearcher()
		cfg, err := cfgSearcher.Get()
		if err != nil {
			log.Fatalf("не удалось получить pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

// GRPCConfig ..
func (s *serviceProvider) GRPCConfig() *config.GRPCConfig {
	if s.grpcConfig == nil {
		cfgSearcher := env.NewGRPCCfgSearcher()
		cfg, err := cfgSearcher.Get()
		if err != nil {
			log.Fatalf("не удалось получить pg config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

// Logger ..
func (s *serviceProvider) Logger() *zerolog.Logger {
	if s.logger == nil {
		cfgSearcher := env.NewLogCfgSearcher()
		cfg, err := cfgSearcher.Get()
		if err != nil {
			log.Fatalf("не удалось получить pg config: %s", err.Error())
		}

		s.logger = setupZeroLog(cfg)
	}

	return s.logger
}

// Pool ..
func (s *serviceProvider) Pool(ctx context.Context) *pgxpool.Pool {
	if s.pool == nil {
		pgxConfig, err := pgxpool.ParseConfig(getPgURL(*s.PgConfig()))
		if err != nil {
			log.Fatalf("ошибка при формировании конфига для pgxpool: %v", err)
		}

		pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
		if err != nil {
			log.Fatalf("ошибка при подключении к DB: %v", err)
		}
		poolCloser := func() error {
			pool.Close()
			return nil
		}

		closer.Add(poolCloser)

		s.pool = pool
	}

	return s.pool
}

// ChatRepository ..
func (s *serviceProvider) ChatRepository(ctx context.Context) repository.ChatRepository {
	if s.chatRepo == nil {
		s.chatRepo = chatRepository.NewRepo(s.Pool(ctx), s.Logger())
	}

	return s.chatRepo
}

// UserRepository ..
func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepo == nil {
		s.userRepo = userRepository.NewRepo(s.Pool(ctx), s.Logger())
	}

	return s.userRepo
}

// MessageRepository ..
func (s *serviceProvider) MessageRepository(ctx context.Context) repository.MessageRepository {
	if s.messageRepo == nil {
		s.messageRepo = messageRepository.NewRepo(s.Pool(ctx), s.Logger())
	}

	return s.messageRepo
}

// ChatService ..
func (s *serviceProvider) ChatService(ctx context.Context) service.ChatService {
	if s.chatService == nil {
		s.chatService = chatService.NewService(
			s.Logger(),
			s.ChatRepository(ctx),
			s.UserRepository(ctx),
			s.MessageRepository(ctx),
		)
	}

	return s.chatService
}

// ChatImpl ..
func (s *serviceProvider) ChatImpl(ctx context.Context) *chat.Implementation {
	if s.chatImpl == nil {
		s.chatImpl = chat.NewImplementation(s.ChatService(ctx))
	}

	return s.chatImpl
}

func getPgURL(cfg config.PgConfig) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DbName,
	)
}

func setupZeroLog(logConfig *config.LogConfig) *zerolog.Logger {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: logConfig.TimeFormat}
	logger := zerolog.New(output).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(logConfig.LogLevel)
	zerolog.TimeFieldFormat = logConfig.TimeFormat

	return &logger
}
