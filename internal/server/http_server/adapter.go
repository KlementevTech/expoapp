package httpserver

import (
	"context"
	"errors"
	"net"
	"net/http"

	"expo/internal/server"
)

type adapter struct {
	s *http.Server
}

func NewAdapter(s *http.Server) server.ServeStopAdapter {
	return &adapter{
		s: s,
	}
}

func (a *adapter) Serve(listener net.Listener) error {
	err := a.s.Serve(listener)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (a *adapter) Stop(ctx context.Context) error {
	return a.s.Shutdown(ctx)
}
