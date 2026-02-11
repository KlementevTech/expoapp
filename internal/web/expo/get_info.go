package expo

import (
	"context"

	expopb "expoapp/pkg/api/expo"
)

func (s *Service) GetInfo(_ context.Context, _ *expopb.GetInfoRequest) (*expopb.GetInfoResponse, error) {
	version := s.versions.GetVersion()

	return &expopb.GetInfoResponse{
		Version: mapGetInfoVersionToPb(version),
	}, nil
}
