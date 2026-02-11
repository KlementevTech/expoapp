package main

import (
	"context"
	"errors"
	"expoapp/internal"
	"expoapp/internal/service"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"strconv"

	"expoapp/internal/web"
	"expoapp/internal/web/expo"

	"github.com/caarlos0/env/v11"
)

var version = "dev"

func main() {
	cfg, err := env.ParseAs[internal.Config]()
	if err != nil {
		slog.Default().Error("failed to parse config", "error", err)
		os.Exit(1)
	}

	expoService := expo.NewService(service.NewVersionService(version))

	ctx, cancel := context.WithCancelCause(context.Background())
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func() {
		select {
		case sig := <-interrupt:
			slog.Default().Info("received signal", "signal", sig)
			cancel(context.Canceled)
		case <-ctx.Done():
		}
	}()

	addr := net.JoinHostPort(cfg.GRPCServerHost, strconv.Itoa(cfg.GRPCServerPort))
	slog.Default().Info("starting grpc server", slog.String("address", addr))
	grpcServer, errChan := web.StartGRPCServer(
		ctx,
		addr,
		expoService,
	)

	go func() {
		if grpcErr := <-errChan; err != nil {
			cancel(fmt.Errorf("failed to start grpc server: %w", grpcErr))
		}
	}()

	<-ctx.Done()
	if err = context.Cause(ctx); err != nil && !errors.Is(err, context.Canceled) {
		slog.Default().Error("something went wrong", "error", err)
		os.Exit(1)
	}

	slog.Default().Info("stopping grpc server")
	grpcServer.GracefulStop()
}
