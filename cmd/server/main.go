package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"

	"expoapp/internal"
	"expoapp/internal/handlers/expo"
	"expoapp/internal/service"

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

	vs := service.NewVersionService(version)

	g, ctx := errgroup.WithContext(context.Background())
	ctx = withInterrupt(ctx, os.Interrupt)

	slog.Default().InfoContext(
		ctx,
		"starting grpc server",
		slog.String("host", cfg.GRPCServerHost),
		slog.Int("port", cfg.GRPCServerPort),
	)
	grpcServer := internal.StartGRPCServer(
		ctx,
		g,
		cfg.GRPCServerHost,
		cfg.GRPCServerPort,
		expo.RegisterHandler(vs),
	)

	g.Go(func() error {
		<-ctx.Done()

		slog.Default().Info("stopping grpc server")
		grpcServer.GracefulStop()

		return nil
	})

	if err = g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		slog.Default().Error("failed to start servers", "error", err)
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
