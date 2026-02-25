package internal

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	GRPCServer GRPCServerConfig `mapstructure:"grpc_server"`
	RESTServer RESTServerConfig `mapstructure:"rest_server"`
	Pprof      PprofConfig      `mapstructure:"pprof"`
	Log        LogConfig        `mapstructure:"log"`
}

type GRPCServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type RESTServerConfig struct {
	Host              string        `mapstructure:"host"`
	Port              int           `mapstructure:"port"`
	ReadHeaderTimeout time.Duration `mapstructure:"read_header_timeout"`
}

type PprofConfig struct {
	Host              string        `mapstructure:"host"`
	Port              int           `mapstructure:"port"`
	Enabled           bool          `mapstructure:"enabled"`
	ReadHeaderTimeout time.Duration `mapstructure:"read_header_timeout"`
}

type LogConfig struct {
	Level string `mapstructure:"level"`
}

func LoadConfig(path string) (Config, error) {
	if path == "" {
		return Config{}, errors.New("no config file path")
	}

	return loadConfigAs[Config](path, "EXPO")
}

func loadConfigAs[T any](path string, envPrefix string) (T, error) {
	viper.AutomaticEnv()
	viper.SetEnvPrefix(envPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AllowEmptyEnv(false)

	viper.SetConfigFile(path)

	var cfg T
	if err := viper.ReadInConfig(); err != nil {
		if tErr, ok := errors.AsType[viper.ConfigFileNotFoundError](err); ok {
			return cfg, fmt.Errorf("config file not found: %w", tErr)
		}
		return cfg, fmt.Errorf("error reading config file: %w", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return cfg, nil
}
