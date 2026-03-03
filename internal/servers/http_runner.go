package servers

import (
	"context"
	"net/http"
	"net/http/pprof"

	"golang.org/x/sync/errgroup"
)

type httpRunner struct {
	name    string
	cfg     HTTPServerConfig
	handler http.Handler
}

func (r *httpRunner) Run(ctx context.Context, g *errgroup.Group) error {
	httpSrv := &http.Server{
		Handler:           r.handler,
		ReadHeaderTimeout: r.cfg.ReadHeaderTimeout,
	}

	return launch(ctx, g, r.cfg.Host, r.cfg.Port, newHTTPAdapter(r.name, httpSrv))
}

func NewRESTRunner(cfg HTTPServerConfig, handler http.Handler) Runner {
	return &httpRunner{
		name:    "REST",
		cfg:     cfg,
		handler: handler,
	}
}

func NewPprofRunner(cfg HTTPServerConfig) Runner {
	mux := http.NewServeMux()

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	return &httpRunner{
		name:    "Pprof",
		cfg:     cfg,
		handler: mux,
	}
}
