package internal

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strconv"

	"expo/internal/rest"

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

	restServer := startRESTServer(
		ctx,
		g,
		rest.SetupRouter(vs),
		cfg.RESTServer,
	)

	g.Go(func() error {
		<-ctx.Done()

		sdCtx, cancel := context.WithTimeout(context.Background(), cfg.RESTServer.ShutdownTimeout)
		defer cancel()

		err := restServer.Shutdown(ctx)
		if err != nil {
			slog.Default().ErrorContext(sdCtx, "failed to shutdown rest server", "error", err)
		}

		grpcServer.GracefulStop()
		slog.Default().Info("stopped gRPC server")

		return nil
	})
}

func startRESTServer(ctx context.Context, g *errgroup.Group, handlers http.Handler, cfg RESTServerConfig) *http.Server {
	s := &http.Server{
		Addr:              net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)),
		Handler:           handlers,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
	}

	g.Go(func() error {
		slog.Default().InfoContext(ctx, fmt.Sprintf("starting REST server at %s", s.Addr))
		err := s.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}

		slog.Default().InfoContext(ctx, "stopped REST server")
		return nil
	})

	return s
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
