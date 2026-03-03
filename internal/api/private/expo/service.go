package expo

import (
	"context"

	expov1 "expo/internal/gen/pb/expo/v1"
)

var _ expov1.ExpoServiceServer = (*serverImpl)(nil)

type serverImpl struct {
	expov1.UnimplementedExpoServiceServer

	versionSvc VersionService
}

func NewExpoServiceServer(versions VersionService) expov1.ExpoServiceServer {
	return &serverImpl{
		versionSvc: versions,
	}
}

func (s *serverImpl) GetVersion(_ context.Context, _ *expov1.GetVersionRequest) (*expov1.GetVersionResponse, error) {
	version := s.versionSvc.GetVersion()

	return &expov1.GetVersionResponse{
		Version: mapVersionToPb(version),
	}, nil
}
