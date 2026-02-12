package expo

import (
	"expoapp/internal/domain"
	"expoapp/pkg/api/pb"

	"google.golang.org/grpc"
)

var _ pb.ExpoServiceServer = (*Service)(nil)

type VersionProvider interface {
	GetVersion() *domain.Version
}

type Service struct {
	pb.UnimplementedExpoServiceServer

	versions VersionProvider
}

func NewService(versions VersionProvider) *Service {
	return &Service{
		versions: versions,
	}
}

func RegisterService(versions VersionProvider) func(s *grpc.Server) {
	return func(s *grpc.Server) {
		pb.RegisterExpoServiceServer(s, NewService(versions))
	}
}
