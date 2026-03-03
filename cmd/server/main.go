package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"expo/internal/api/grpc/expo"
	"expo/internal/api/rest"
	"expo/internal/servers"
	"expo/internal/service"

	"expo/internal"

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
	internal.ProvideDeps(i, version)

	defer func() {
		internal.ShutdownDeps(i)
	}()

	versionSvc := do.MustInvoke[*service.VersionService](i)

	runners := []servers.Runner{
		servers.NewRESTRunner(cfg.RESTServer, rest.NewHTTPHandler(versionSvc)),
		servers.NewGRPCRunner(cfg.GRPCServer, expo.NewServerRegister(versionSvc)),
	}

	if cfg.PprofEnabled {
		runners = append(runners, servers.NewPprofRunner(cfg.PprofServer))
	}

	err = servers.Run(ctx, runners...)
	if err != nil {
		return err
	}

	slog.Default().Info("server stopped")
	return nil
}
