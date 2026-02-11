package service

import "expoapp/internal/domain"

type VersionService struct {
	model *domain.Version
}

func NewVersionService(version string) *VersionService {
	return &VersionService{
		model: domain.NewVersion(version),
	}
}

func (s *VersionService) GetVersion() *domain.Version {
	return s.model
}
