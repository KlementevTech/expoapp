package expo

import (
	"context"

	expov1 "expo/pkg/pb/expo/v1"

	"google.golang.org/grpc"
)

var _ expov1.ExpoServiceServer = (*service)(nil)

type service struct {
	expov1.UnimplementedExpoServiceServer

	vs VersionService
}

func NewService(versions VersionService) expov1.ExpoServiceServer {
	return &service{
		vs: versions,
	}
}

func (s *service) GetVersion(_ context.Context, _ *expov1.GetVersionRequest) (*expov1.GetVersionResponse, error) {
	version := s.vs.GetVersion()

	return &expov1.GetVersionResponse{
		Version: mapVersionToPb(version),
	}, nil
}

func Register(versionService VersionService) func(s *grpc.Server) {
	return func(s *grpc.Server) {
		expov1.RegisterExpoServiceServer(s, NewService(versionService))
	}
}
