package internal

import (
	"context"
	"net"
	"net/http"
)

type httpAdapter struct {
	name string
	s    *http.Server
}

func newHTTPAdapter(name string, s *http.Server) *httpAdapter {
	return &httpAdapter{
		name: name,
		s:    s,
	}
}

func (a *httpAdapter) Name() string {
	return a.name
}

func (a *httpAdapter) Serve(listener net.Listener) error {
	return a.s.Serve(listener)
}

func (a *httpAdapter) Stop(ctx context.Context) error {
	return a.s.Shutdown(ctx)
}
