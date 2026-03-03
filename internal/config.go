package internal

import (
	"errors"
	"fmt"
	"strings"

	"expo/internal/servers"

	"github.com/spf13/viper"
)

type Config struct {
	GRPCServer   servers.GRPCServerConfig `mapstructure:"grpc_server"`
	RESTServer   servers.HTTPServerConfig `mapstructure:"rest_server"`
	PprofServer  servers.HTTPServerConfig `mapstructure:"pprof_server"`
	Log          LogConfig                `mapstructure:"log"`
	PprofEnabled bool                     `mapstructure:"pprof_enabled"`
}

type LogConfig struct {
	Level string `mapstructure:"level"`
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
