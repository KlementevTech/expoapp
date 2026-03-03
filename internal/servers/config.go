package servers

import "time"

type HTTPServerConfig struct {
	Host              string
	Port              int
	ReadHeaderTimeout time.Duration
}

type GRPCServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}
