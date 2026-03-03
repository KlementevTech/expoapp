package servers

import (
	"context"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCRegistrator interface {
	RegisterServer(s *grpc.Server)
}

type grpcRunner struct {
	name      string
	cfg       GRPCServerConfig
	registers []GRPCRegistrator
}

func NewGRPCRunner(cfg GRPCServerConfig, register ...GRPCRegistrator) Runner {
	return &grpcRunner{
		name:      "gRPC",
		cfg:       cfg,
		registers: append([]GRPCRegistrator{}, register...),
	}
}

func (r *grpcRunner) Run(ctx context.Context, g *errgroup.Group) error {
	grpcSrv := grpc.NewServer()

	for _, reg := range r.registers {
		reg.RegisterServer(grpcSrv)
	}

	reflection.Register(grpcSrv)

	return launch(ctx, g, r.cfg.Host, r.cfg.Port, newGRPCAdapter(r.name, grpcSrv))
}
