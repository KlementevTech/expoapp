package expo

import (
	"context"

	"expoapp/internal/model"
	expov1 "expoapp/pkg/pb/expo/v1"

	"google.golang.org/grpc"
)

type VersionService interface {
	GetCurrentVersion() *model.Version
}

var _ expov1.ExpoServiceServer = (*service)(nil)

type service struct {
	expov1.UnimplementedExpoServiceServer

	versions VersionService
}

func NewService(versions VersionService) expov1.ExpoServiceServer {
	return &service{
		versions: versions,
	}
}

func (s *service) GetVersion(_ context.Context, _ *expov1.GetVersionRequest) (*expov1.GetVersionResponse, error) {
	version := s.versions.GetCurrentVersion()

	return &expov1.GetVersionResponse{
		Version: mapVersionToPb(version),
	}, nil
}

func RegisterHandler(versionService VersionService) func(s *grpc.Server) {
	return func(s *grpc.Server) {
		expov1.RegisterExpoServiceServer(s, NewService(versionService))
	}
}
