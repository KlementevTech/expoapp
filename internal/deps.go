package internal

import (
	"log/slog"

	"expo/internal/service"

	"github.com/samber/do/v2"
)

func ProvideDeps(i do.Injector, version string) do.Injector {
	do.Provide(i, func(_ do.Injector) (*service.VersionService, error) {
		return service.NewVersionService(version), nil
	})

	return i
}

func ShutdownDeps(i do.Injector) {
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
