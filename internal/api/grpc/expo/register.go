package expo

import (
	expov1 "expo/internal/gen/pb/expo/v1"
	"expo/internal/servers"

	"google.golang.org/grpc"
)

var _ servers.GRPCRegistrator = (*serverRegister)(nil)

type serverRegister struct {
	versionSvc VersionService
}

func NewServerRegister(versionSvc VersionService) servers.GRPCRegistrator {
	return &serverRegister{
		versionSvc: versionSvc,
	}
}

func (r *serverRegister) RegisterServer(srv *grpc.Server) {
	expov1.RegisterExpoServiceServer(srv, NewExpoServiceServer(r.versionSvc))
}
