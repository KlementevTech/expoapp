package service

import (
	"expo/internal/model"
)

type VersionService struct {
	currentVersion string
}

func NewVersionService(version string) *VersionService {
	return &VersionService{
		currentVersion: version,
	}
}

func (s *VersionService) GetVersion() *model.Version {
	return new(model.Version(s.currentVersion))
}
