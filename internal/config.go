package internal

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	GRPCServer GRPCServerConfig `mapstructure:"grpc_server"`
}

type GRPCServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

func LoadConfig[T any](path string) (*T, error) {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("EXPO")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
			slog.Default().Warn("config file not found, using defaults")
		} else {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	func() {
		viper.SetDefault("grpc_server.host", "127.0.0.1")
		viper.SetDefault("grpc_server.port", "50051")
	}()

	cfg := new(T)
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return cfg, nil
}
