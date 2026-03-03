package expo

import (
	"expo/internal/model"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/version_service_mock.go -package=mocks

type VersionService interface {
	GetVersion() *model.Version
}
