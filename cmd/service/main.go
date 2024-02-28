package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"

	chatAPI "github.com/mistandok/chat-server/internal/api/chat"
	chatRepository "github.com/mistandok/chat-server/internal/repository/chat"
	chatService "github.com/mistandok/chat-server/internal/service/chat"

	"github.com/mistandok/chat-server/internal/config"
	"github.com/mistandok/chat-server/internal/config/env"
	"github.com/mistandok/chat-server/internal/repository/postgresql"
	"github.com/mistandok/chat-server/pkg/chat_v1"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", ".env", "path to config file")
	flag.Parse()
}

func main() {
	ctx := context.Background()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("ошибка при получении переменных окружения: %v", err)
	}

	grpcConfig, err := env.NewGRPCCfgSearcher().Get()
	if err != nil {
		log.Fatalf("ошибка при получении конфига grpc: %v", err)
	}

	logConfig, err := env.NewLogCfgSearcher().Get()
	if err != nil {
		log.Fatalf("ошибка при получении уровня логирования: %v", err)
	}

	pgConfig, err := env.NewPgCfgSearcher().Get()
	if err != nil {
		log.Fatalf("ошибка при получении конфига DB: %v", err)
	}

	pool, connCloser := postgresql.MustInitPgConnection(ctx, *pgConfig)
	defer connCloser()

	listenConfig := net.ListenConfig{}
	listener, err := listenConfig.Listen(ctx, "tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("ошибка при прослушивании: %v", err)
	}

	logger := setupZeroLog(logConfig)

	chatRepo := chatRepository.NewRepo(pool, logger)
	chatServ := chatService.NewService(chatRepo, logger)
	chatAPIServer := chatAPI.NewImplementation(chatServ)

	server := grpc.NewServer()
	reflection.Register(server)
	chat_v1.RegisterChatV1Server(server, chatAPIServer)

	log.Printf("сервер запущен на %v", listener.Addr())

	if err := server.Serve(listener); err != nil {
		log.Fatalf("ошибка сервера: %v", err)
	}
}

func setupZeroLog(logConfig *config.LogConfig) *zerolog.Logger {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: logConfig.TimeFormat}
	logger := zerolog.New(output).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(logConfig.LogLevel)
	zerolog.TimeFieldFormat = logConfig.TimeFormat

	return &logger
}
