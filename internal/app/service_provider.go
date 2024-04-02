package app

import (
	"context"
	"github.com/mistandok/chat-server/internal/client"
	"github.com/mistandok/chat-server/internal/client/access"
	"github.com/mistandok/chat-server/internal/interceptor"
	"github.com/mistandok/chat-server/pkg/auth_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"

	"github.com/mistandok/platform_common/pkg/closer"
	"github.com/mistandok/platform_common/pkg/db"
	"github.com/mistandok/platform_common/pkg/db/pg"

	"github.com/mistandok/chat-server/internal/api/chat"
	"github.com/mistandok/chat-server/internal/config"
	"github.com/mistandok/chat-server/internal/config/env"
	"github.com/mistandok/chat-server/internal/repository"
	chatRepository "github.com/mistandok/chat-server/internal/repository/chat"
	messageRepository "github.com/mistandok/chat-server/internal/repository/message"
	userRepository "github.com/mistandok/chat-server/internal/repository/user"
	"github.com/mistandok/chat-server/internal/service"
	chatService "github.com/mistandok/chat-server/internal/service/chat"
	"github.com/rs/zerolog"
)

type serviceProvider struct {
	pgConfig               *config.PgConfig
	grpcConfig             *config.GRPCConfig
	httpConfig             *config.HTTPConfig
	swaggerConfig          *config.SwaggerConfig
	authConfig             *config.AuthConfig
	logger                 *zerolog.Logger
	accessClient           auth_v1.AccessV1Client
	accessClientFacade     client.AccessClient
	accessCheckInterceptor *interceptor.AccessCheckInterceptor

	dbClient  db.Client
	txManager db.TxManager

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
			log.Fatalf("не удалось получить grpc config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

// HTTPConfig ..
func (s *serviceProvider) HTTPConfig() *config.HTTPConfig {
	if s.httpConfig == nil {
		cfgSearcher := env.NewHTTPCfgSearcher()
		cfg, err := cfgSearcher.Get()
		if err != nil {
			log.Fatalf("не удалось получить http config: %s", err.Error())
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}

// SwaggerConfig ..
func (s *serviceProvider) SwaggerConfig() *config.SwaggerConfig {
	if s.swaggerConfig == nil {
		cfgSearcher := env.NewSwaggerConfigSearcher()
		cfg, err := cfgSearcher.Get()
		if err != nil {
			log.Fatalf("не удалось получить swagger config: %s", err.Error())
		}

		s.swaggerConfig = cfg
	}

	return s.swaggerConfig
}

// AuthConfig ..
func (s *serviceProvider) AuthConfig() *config.AuthConfig {
	if s.authConfig == nil {
		cfgSearcher := env.NewAuthCfgSearcher()
		cfg, err := cfgSearcher.Get()
		if err != nil {
			log.Fatalf("не удалось получить auth config: %s", err.Error())
		}

		s.authConfig = cfg
	}

	return s.authConfig
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

// DBClient ..
func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PgConfig().DSN(), s.Logger())
		if err != nil {
			log.Fatalf("ошибка при создании клиента DB: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("нет связи с БД: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

// TxManager ..
func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = pg.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

// ChatRepository ..
func (s *serviceProvider) ChatRepository(ctx context.Context) repository.ChatRepository {
	if s.chatRepo == nil {
		s.chatRepo = chatRepository.NewRepo(s.Logger(), s.DBClient(ctx))
	}

	return s.chatRepo
}

// UserRepository ..
func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepo == nil {
		s.userRepo = userRepository.NewRepo(s.Logger(), s.DBClient(ctx))
	}

	return s.userRepo
}

// MessageRepository ..
func (s *serviceProvider) MessageRepository(ctx context.Context) repository.MessageRepository {
	if s.messageRepo == nil {
		s.messageRepo = messageRepository.NewRepo(s.Logger(), s.DBClient(ctx))
	}

	return s.messageRepo
}

// ChatService ..
func (s *serviceProvider) ChatService(ctx context.Context) service.ChatService {
	if s.chatService == nil {
		s.chatService = chatService.NewService(
			s.Logger(),
			s.TxManager(ctx),
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

func (s *serviceProvider) AccessV1Client(_ context.Context) auth_v1.AccessV1Client {
	if s.accessClient == nil {
		cfg := s.AuthConfig()
		conn, err := grpc.Dial(
			cfg.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			log.Fatalf("ошибка при установлении соединения с auth-сервисом: %v", err)
		}

		s.accessClient = auth_v1.NewAccessV1Client(conn)
	}

	return s.accessClient
}

func (s *serviceProvider) AccessClientFacade(ctx context.Context) client.AccessClient {
	if s.accessClientFacade == nil {
		s.accessClientFacade = access.NewClientFacade(s.Logger(), s.AccessV1Client(ctx))
	}

	return s.accessClientFacade
}

func (s *serviceProvider) AccessCheckInterceptor(ctx context.Context) *interceptor.AccessCheckInterceptor {
	if s.accessCheckInterceptor == nil {
		s.accessCheckInterceptor = interceptor.NewAccessCheckInterceptor(s.Logger(), s.AccessClientFacade(ctx))
	}

	return s.accessCheckInterceptor
}

func setupZeroLog(logConfig *config.LogConfig) *zerolog.Logger {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: logConfig.TimeFormat}
	logger := zerolog.New(output).With().Timestamp().Logger()
	logger = logger.Level(logConfig.LogLevel)
	zerolog.TimeFieldFormat = logConfig.TimeFormat

	return &logger
}
