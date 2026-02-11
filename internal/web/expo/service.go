package expo

import (
	"expoapp/internal/domain"
	expopb "expoapp/pkg/api/expo"
)

var _ expopb.ExpoServiceServer = (*Service)(nil)

type VersionService interface {
	GetVersion() *domain.Version
}

type Service struct {
	expopb.UnimplementedExpoServiceServer

	versions VersionService
}

func NewService(versions VersionService) *Service {
	return &Service{
		versions: versions,
	}
}
