package expo

import (
	"expoapp/internal/domain"
	expopb "expoapp/pkg/api/expo"
)

func mapGetInfoVersionToPb(version *domain.Version) *expopb.GetInfoResponse_Version {
	return &expopb.GetInfoResponse_Version{
		Version: version.Version(),
	}
}
