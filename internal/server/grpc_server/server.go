package grpcserver

import (
	"expo/internal/server"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Registrator interface {
	RegisterServer(s *grpc.Server)
}

type Config struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type srv struct {
	server.Server
}

func NewServer(cfg Config, register ...Registrator) server.Runner {
	grpcSrv := grpc.NewServer()
	reflection.Register(grpcSrv)

	for _, r := range register {
		r.RegisterServer(grpcSrv)
	}

	return &srv{
		Server: server.NewServer("gRPC", cfg.Host, cfg.Port, NewAdapter(grpcSrv)),
	}
}
