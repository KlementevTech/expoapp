package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"expo/internal"
	"expo/internal/service"
)

var version = "dev"

func main() {
	cfg, err := internal.LoadConfig()
	if err != nil {
		slog.Default().Error("failed to load config", "error", err)
		os.Exit(1)
	}

	err = setupLogger(cfg.Log, version)
	if err != nil {
		slog.Default().Error("failed to setup logger", "error", err)
		os.Exit(1)
	}

	vs := service.NewVersionService(version)

	ctx := notifyContext(os.Interrupt)
	err = internal.RunServers(ctx, cfg, vs)
	if err != nil {
		slog.Default().Error("failed to start servers", "error", err)
		os.Exit(1)
	}
}

func notifyContext(sig ...os.Signal) context.Context {
	interrupt := make(chan os.Signal, len(sig))
	signal.Notify(interrupt, sig...)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		s := <-interrupt
		slog.Default().Info("received signal", slog.String("signal", s.String()))
		cancel()
	}()

	return ctx
}

func setupLogger(cfg internal.LogConfig, version string) error {
	var l slog.Level
	if err := l.UnmarshalText([]byte(cfg.Level)); err != nil {
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
