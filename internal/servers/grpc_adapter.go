package servers

import (
	"context"
	"errors"
	"net"

	"google.golang.org/grpc"
)

type grpcAdapter struct {
	name string
	s    *grpc.Server
}

func newGRPCAdapter(name string, s *grpc.Server) *grpcAdapter {
	return &grpcAdapter{
		name: name,
		s:    s,
	}
}

func (a *grpcAdapter) Name() string {
	return a.name
}

func (a *grpcAdapter) Serve(listener net.Listener) error {
	return a.s.Serve(listener)
}

func (a *grpcAdapter) Stop(ctx context.Context) error {
	done := make(chan struct{})

	go func() {
		a.s.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		a.s.Stop()
		return errors.New("gRPC shutdown timeout, forcing stop")
	}
}
