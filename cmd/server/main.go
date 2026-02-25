package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"os"
	"os/signal"

	"expo/internal"
	"expo/internal/service"

	"golang.org/x/sync/errgroup"
)

var version = "dev"

func main() {
	var path string
	flag.StringVar(&path, "config", "", "config file path")
	flag.Parse()

	cfg, err := internal.LoadConfig(path)
	if err != nil {
		slog.Default().Error("failed to load config", "error", err)
		os.Exit(1)
	}

	err = setupLogger(cfg.LogLevel, version)
	if err != nil {
		slog.Default().Error("failed to setup logger", "error", err)
		os.Exit(1)
	}

	vs := service.NewVersionService(version)

	g, ctx := errgroup.WithContext(context.Background())
	ctx = withInterrupt(ctx, os.Interrupt)

	internal.StartServers(ctx, g, cfg, vs)

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

func setupLogger(level, version string) error {
	var l slog.Level
	if err := l.UnmarshalText([]byte(level)); err != nil {
		return err
	}

	slog.SetDefault(
		slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: l,
			}),
		).With(
			slog.String("version", version),
		),
	)
	return nil
}
