package expo

import (
	"expoapp/internal/domain"
	"expoapp/pkg/api/pb"
)

func mapGetInfoV1VersionToPb(version *domain.Version) *pb.GetInfoV1Response_Version {
	return &pb.GetInfoV1Response_Version{
		Version: version.Version(),
	}
}
