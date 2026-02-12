package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"

	"expoapp/internal"
	"expoapp/internal/service"
	"expoapp/internal/web"
	"expoapp/internal/web/expo"

	"github.com/caarlos0/env/v11"
	"golang.org/x/sync/errgroup"
)

var version = "dev"

func main() {
	cfg, err := env.ParseAs[internal.Config]()
	if err != nil {
		slog.Default().Error("failed to parse config", "error", err)
		os.Exit(1)
	}

	versions := service.NewVersionService(version)

	grpcServer := web.NewGRPCServer(
		web.ServiceRegisterFunc(expo.RegisterService(versions)),
	)

	g, ctx := errgroup.WithContext(context.Background())
	ctx = withInterrupt(ctx, os.Interrupt)

	slog.Default().Info(
		"starting grpc server",
		slog.String("host", cfg.GRPCServerHost),
		slog.Int("port", cfg.GRPCServerPort),
	)
	web.StartGRPCServer(
		ctx,
		g,
		grpcServer,
		cfg.GRPCServerHost,
		cfg.GRPCServerPort,
	)

	g.Go(func() error {
		<-ctx.Done()

		slog.Default().Info("stopping grpc server")
		grpcServer.GracefulStop()

		return nil
	})

	if err = g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		slog.Default().Error("something went wrong", "error", err)
		os.Exit(1)
	}
}

func withInterrupt(ctx context.Context, sig ...os.Signal) context.Context {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, sig...)

	ctx, cancel := context.WithCancel(ctx)

	go func() {
		select {
		case <-ctx.Done():
		case s := <-interrupt:
			slog.Default().Info("received signal", "signal", s)
			cancel()
		}
	}()

	return ctx
}
