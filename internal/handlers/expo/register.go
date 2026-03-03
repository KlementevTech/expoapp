package expo

import (
	"expo/internal"
	expov1 "expo/internal/gen/pb/expo/v1"

	"google.golang.org/grpc"
)

var _ internal.GRPCRegistrator = (*register)(nil)

type register struct {
	versionSvc VersionService
}

func NewRegister(versionSvc VersionService) internal.GRPCRegistrator {
	return &register{
		versionSvc: versionSvc,
	}
}

func (r *register) RegisterServer(srv *grpc.Server) {
	expov1.RegisterExpoServiceServer(srv, NewService(r.versionSvc))
}
