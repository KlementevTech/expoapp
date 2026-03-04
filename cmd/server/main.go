package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"expo/internal/api/rest"
	httpserver "expo/internal/server/http_server"

	"expo/internal"
	apigrpc "expo/internal/api/grpc"
	"expo/internal/api/handler"
	"expo/internal/server"
	grpcconnect "expo/internal/server/grpc_connect"
	grpcserver "expo/internal/server/grpc_server"
	"expo/internal/server/profiling"

	"github.com/samber/do/v2"
)

var version = "dev"

func main() {
	if err := run(); err != nil {
		slog.Default().Error("something went wrong", "error", err)
		os.Exit(1)
	}
}

func run() error {
	var path string
	flag.StringVar(&path, "config", "", "config file path")
	flag.Parse()

	cfg, err := internal.LoadConfig(path)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	err = internal.SetupLogger(cfg.Log, version)
	if err != nil {
		return fmt.Errorf("setup logger: %w", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	i := do.New()
	shutdown := internal.ProvideDeps(i, cfg)
	defer func() {
		shutdown()
	}()

	parts := do.MustInvoke[handler.PartRepository](i)

	servers := []server.Runner{
		grpcconnect.NewServer(cfg.ConnectServer, handler.NewHandler(parts, true)),
		grpcserver.NewServer(cfg.GRPCServer, apigrpc.NewRegister(parts)),
		httpserver.NewServer(cfg.RESTServer, rest.NewHandler()),
	}

	if cfg.PprofEnabled {
		servers = append(servers, profiling.NewServer(cfg.PprofServer))
	}

	err = server.Run(ctx, servers...)
	if err != nil {
		return err
	}

	slog.Default().Info("server stopped")
	return nil
}
