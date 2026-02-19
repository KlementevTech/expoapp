package internal

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"strconv"

	"expo/internal/handlers/expo"
	"expo/internal/service"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func StartServers(ctx context.Context, g *errgroup.Group, cfg Config, vs *service.VersionService) {
	grpcServer := startGRPCServer(
		ctx,
		g,
		cfg.GRPCServer,
		expo.RegisterHandler(vs),
	)

	g.Go(func() error {
		<-ctx.Done()
		slog.Default().Info("stopping gRPC server")
		grpcServer.GracefulStop()

		return nil
	})
}

func startGRPCServer(
	ctx context.Context,
	g *errgroup.Group,
	cfg GRPCServerConfig,
	register ...func(s *grpc.Server),
) *grpc.Server {
	s := grpc.NewServer()
	reflection.Register(s)

	for _, r := range register {
		r(s)
	}

	g.Go(func() error {
		lis, err := new(net.ListenConfig).
			Listen(ctx, "tcp", net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)))
		if err != nil {
			return err
		}

		defer func() {
			_ = lis.Close()
		}()

		slog.Default().InfoContext(ctx, fmt.Sprintf("starting gRPC server at %s", lis.Addr().String()))
		if err = s.Serve(lis); err != nil {
			return err
		}

		return nil
	})

	return s
}
