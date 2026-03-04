package grpcserver

import (
	"context"
	"errors"
	"net"

	"expo/internal/server"

	"google.golang.org/grpc"
)

type adapter struct {
	s *grpc.Server
}

func NewAdapter(s *grpc.Server) server.ServeStopAdapter {
	return &adapter{
		s: s,
	}
}

func (a *adapter) Serve(listener net.Listener) error {
	return a.s.Serve(listener)
}

func (a *adapter) Stop(ctx context.Context) error {
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
