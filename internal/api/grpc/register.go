package grpc

import (
	expov1 "expo/gen/pb/part/v1"
	grpcserver "expo/internal/server/grpc_server"

	"google.golang.org/grpc"
)

var _ grpcserver.Registrator = (*register)(nil)

type register struct {
	parts PartRepository
}

func NewRegister(parts PartRepository) grpcserver.Registrator {
	return &register{
		parts: parts,
	}
}

func (r *register) RegisterServer(srv *grpc.Server) {
	expov1.RegisterPartServiceServer(srv, NewExpoServiceServer(r.parts))
}
