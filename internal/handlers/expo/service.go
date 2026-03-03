package expo

import (
	"context"

	expov1 "expo/internal/gen/pb/expo/v1"
)

var _ expov1.ExpoServiceServer = (*serviceImpl)(nil)

type serviceImpl struct {
	expov1.UnimplementedExpoServiceServer

	versionSvc VersionService
}

func NewService(versions VersionService) expov1.ExpoServiceServer {
	return &serviceImpl{
		versionSvc: versions,
	}
}

func (s *serviceImpl) GetVersion(_ context.Context, _ *expov1.GetVersionRequest) (*expov1.GetVersionResponse, error) {
	version := s.versionSvc.GetVersion()

	return &expov1.GetVersionResponse{
		Version: mapVersionToPb(version),
	}, nil
}
