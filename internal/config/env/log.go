package env

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/mistandok/chat-server/internal/config"
	"github.com/rs/zerolog"
)

const (
	logLevel      = "LOG_LEVEL"
	logTimeFormat = "LOG_TIME_FORMAT"
)

// LogCfgSearcher logger config searcher.
type LogCfgSearcher struct{}

// NewLogCfgSearcher get instance for logger config searcher.
func NewLogCfgSearcher() *LogCfgSearcher {
	return &LogCfgSearcher{}
}

// Get config for logger.
func (s *LogCfgSearcher) Get() (*config.LogConfig, error) {
	level := os.Getenv(logLevel)
	if len(level) == 0 {
		return nil, errors.New("не найден уровень логирования")
	}

	logLevelInt, err := strconv.Atoi(level)
	if err != nil {
		return nil, fmt.Errorf("некорректное значение уровня логирования: %w", err)
	}

	timeFormat := os.Getenv(logTimeFormat)
	if len(level) == 0 {
		return nil, errors.New("не найден формат времени логирования")
	}

	return &config.LogConfig{
		LogLevel:   zerolog.Level(logLevelInt),
		TimeFormat: timeFormat,
	}, nil
}
