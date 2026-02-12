package expo

import (
	"context"

	"expoapp/pkg/api/pb"
)

func (s *Service) GetInfoV1(_ context.Context, _ *pb.GetInfoV1Request) (*pb.GetInfoV1Response, error) {
	version := s.versions.GetVersion()

	return &pb.GetInfoV1Response{
		Version: mapGetInfoV1VersionToPb(version),
	}, nil
}
