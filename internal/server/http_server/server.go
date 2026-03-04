package httpserver

import (
	"net/http"
	"time"

	"expo/internal/server"
)

type Config struct {
	Host              string
	Port              int
	ReadHeaderTimeout time.Duration
}

type srv struct {
	server.Server
}

func NewServer(cfg Config, handler http.Handler) server.Runner {
	httpSrv := &http.Server{
		Handler:           handler,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
	}

	return &srv{
		Server: server.NewServer("HTTP", cfg.Host, cfg.Port, NewAdapter(httpSrv)),
	}
}
