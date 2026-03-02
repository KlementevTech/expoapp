package internal

import (
	"fmt"
	"log/slog"
	"os"
)

func SetupLogger(cfg LogConfig, version string) error {
	var l slog.Level
	if err := l.UnmarshalText([]byte(cfg.Level)); err != nil {
		return fmt.Errorf("invalid log level: %w", err)
	}

	options := &slog.HandlerOptions{
		Level: l,
	}

	handler := slog.NewJSONHandler(os.Stdout, options)
	logger := slog.New(handler).With(
		slog.String("version", version),
	)

	slog.SetDefault(logger)
	return nil
}
