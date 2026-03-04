package internal

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"expo/internal/api/handler"
	"expo/internal/storage/sqlite"
	"expo/internal/storage/sqlite/sqlc"

	"github.com/samber/do/v2"
)

type ShutdownFunc func()

func ProvideDeps(i do.Injector, cfg *Config) ShutdownFunc {
	do.Provide[*sql.DB](i, func(_ do.Injector) (*sql.DB, error) {
		sCfg := cfg.SQLite

		// 1. Работа с путями
		absPath, err := filepath.Abs(sCfg.Path)
		if err != nil {
			return nil, fmt.Errorf("failed to get absolute path: %w", err)
		}

		dir := filepath.Dir(absPath)
		if err = os.MkdirAll(dir, 0o750); err != nil {
			return nil, fmt.Errorf("failed to create db directory %s: %w", dir, err)
		}

		// 2. Открытие соединения
		dsn := fmt.Sprintf("%s?_journal=WAL&_foreign_keys=on", absPath)
		db, err := sql.Open("sqlite", dsn)
		if err != nil {
			return nil, fmt.Errorf("invalid dsn: %w", err)
		}

		// 3. Настройка пула из конфига
		db.SetMaxOpenConns(sCfg.MaxOpenConns)
		db.SetMaxIdleConns(sCfg.MaxIdleConns)
		db.SetConnMaxLifetime(sCfg.ConnMaxLifetime)

		// 4. Проверка связи с настраиваемым таймаутом
		ctx, cancel := context.WithTimeout(context.Background(), sCfg.ConnectTimeout)
		defer cancel()

		if err = db.PingContext(ctx); err != nil {
			return nil, fmt.Errorf("database unreachable (timeout %v): %w", sCfg.ConnectTimeout, err)
		}

		return db, nil
	})

	do.Provide[handler.PartRepository](i, func(i do.Injector) (handler.PartRepository, error) {
		db := do.MustInvoke[*sql.DB](i)

		queries := sqlc.New(db)

		return sqlite.NewPartStorage(db, queries), nil
	})

	return func() {
		if report := i.Shutdown(); !report.Succeed {
			for ref, err := range report.Errors {
				slog.Default().Error(
					"injector shutdown",
					slog.String("service", ref.Service),
					slog.Any("error", err),
				)
			}
		}
	}
}
