package config

import (
	"fmt"
	"net"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

// GRPCConfigSearcher interface for search grpc config
type GRPCConfigSearcher interface {
	Get() (*GRPCConfig, error)
}

// LogConfigSearcher interface for serach Log config.
type LogConfigSearcher interface {
	Get() (*LogConfig, error)
}

// PgConfigSearcher interface for search PG config.
type PgConfigSearcher interface {
	Get() (*PgConfig, error)
}

// HTTPConfigSearcher interface for search Http config.
type HTTPConfigSearcher interface {
	Get() (*HTTPConfig, error)
}

// SwaggerConfigSearcher interface for search Http config.
type SwaggerConfigSearcher interface {
	Get() (*SwaggerConfig, error)
}

// AuthConfigSearcher interface for search grpc config
type AuthConfigSearcher interface {
	Get() (*GRPCConfig, error)
}

// Load dotenv from path to env
func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}

// GRPCConfig grpc config.
type GRPCConfig struct {
	Host string
	Port string
}

// Address get address for grpc server.
func (cfg *GRPCConfig) Address() string {
	return net.JoinHostPort(cfg.Host, cfg.Port)
}

// LogConfig config for zerolog.
type LogConfig struct {
	LogLevel   zerolog.Level
	TimeFormat string
}

// PgConfig config for postgresql.
type PgConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DbName   string
}

// DSN ..
func (cfg *PgConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DbName,
	)
}

// HTTPConfig config for HTTP
type HTTPConfig struct {
	Host string
	Port string
}

// Address get address from config
func (cfg *HTTPConfig) Address() string {
	return net.JoinHostPort(cfg.Host, cfg.Port)
}

// SwaggerConfig config for Swagger
type SwaggerConfig struct {
	Host string
	Port string
}

// Address get address from config
func (cfg *SwaggerConfig) Address() string {
	return net.JoinHostPort(cfg.Host, cfg.Port)
}

// AuthConfig grpc config.
type AuthConfig struct {
	Host string
	Port string
}

// Address get address from config
func (cfg *AuthConfig) Address() string {
	return net.JoinHostPort(cfg.Host, cfg.Port)
}

// PrometheusConfig config for Prometheus
type PrometheusConfig struct {
	PublicHost string
	PublicPort string
	Host       string
	Port       string
}

// Address get address from config
func (cfg *PrometheusConfig) Address() string {
	return net.JoinHostPort(cfg.Host, cfg.Port)
}

// PublicAddress get address from config
func (cfg *PrometheusConfig) PublicAddress() string {
	return net.JoinHostPort(cfg.PublicHost, cfg.PublicPort)
}
