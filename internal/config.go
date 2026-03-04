package internal

import (
	"errors"
	"fmt"
	"strings"
	"time"

	grpcserver "expo/internal/server/grpc_server"
	httpserver "expo/internal/server/http_server"

	"github.com/spf13/viper"
)

type Config struct {
	GRPCServer    grpcserver.Config `mapstructure:"grpc_server"`
	ConnectServer httpserver.Config `mapstructure:"connect_server"`
	RESTServer    httpserver.Config `mapstructure:"rest_server"`
	PprofServer   httpserver.Config `mapstructure:"pprof_server"`
	SQLite        SQLiteConfig      `mapstructure:"sqlite"`
	Log           LogConfig         `mapstructure:"log"`
	PprofEnabled  bool              `mapstructure:"pprof_enabled"`
}

type LogConfig struct {
	Level string `mapstructure:"level"`
}

type SQLiteConfig struct {
	Path            string        `mapstructure:"path"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	ConnectTimeout  time.Duration `mapstructure:"connect_timeout"`
}

func LoadConfig(path string) (*Config, error) {
	if path == "" {
		return nil, errors.New("no config file path")
	}
	return loadConfigAs[Config](path, "EXPO")
}

func loadConfigAs[T any](path string, envPrefix string) (*T, error) {
	viper.AutomaticEnv()
	viper.SetEnvPrefix(envPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AllowEmptyEnv(false)

	viper.SetConfigFile(path)

	var cfg T
	if err := viper.ReadInConfig(); err != nil {
		if tErr, ok := errors.AsType[viper.ConfigFileNotFoundError](err); ok {
			return nil, fmt.Errorf("config file not found: %w", tErr)
		}
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return &cfg, nil
}
