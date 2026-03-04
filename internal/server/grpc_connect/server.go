package grpcconnect

import (
	"net/http"

	"expo/internal/server"
	httpserver "expo/internal/server/http_server"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type srv struct {
	server.Server
}

func NewServer(cfg httpserver.Config, handler http.Handler) server.Runner {
	httpSrv := &http.Server{
		Handler:           h2c.NewHandler(handler, &http2.Server{}),
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
	}

	return &srv{
		Server: server.NewServer("gRPC connect", cfg.Host, cfg.Port, httpserver.NewAdapter(httpSrv)),
	}
}
