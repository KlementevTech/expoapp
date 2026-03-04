package profiling

import (
	"net/http"
	"net/http/pprof"

	"expo/internal/server"

	httpserver "expo/internal/server/http_server"
)

type srv struct {
	server.Server
}

func NewServer(cfg httpserver.Config) server.Runner {
	handler := http.NewServeMux()

	handler.HandleFunc("/debug/pprof/", pprof.Index)
	handler.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	handler.HandleFunc("/debug/pprof/profile", pprof.Profile)
	handler.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	handler.HandleFunc("/debug/pprof/trace", pprof.Trace)

	httpSrv := &http.Server{
		Handler:           handler,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
	}

	return &srv{
		Server: server.NewServer("Pprof", cfg.Host, cfg.Port, httpserver.NewAdapter(httpSrv)),
	}
}
