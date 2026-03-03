package internal

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/http/pprof"
	"strconv"
	"time"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const shutdownTimeout = 5 * time.Second

type serverRunner interface {
	Name() string
	Serve(listener net.Listener) error
	Stop(ctx context.Context) error
}

func run(
	ctx context.Context,
	g *errgroup.Group,
	host string,
	port int,
	s serverRunner,
) error {
	addr := net.JoinHostPort(host, strconv.Itoa(port))

	var lc net.ListenConfig
	lis, err := lc.Listen(ctx, "tcp4", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	slog.Default().InfoContext(ctx, fmt.Sprintf("%s server is listening", s.Name()), slog.String("addr", addr))

	g.Go(func() error {
		errChan := make(chan error, 1)
		go func() {
			errChan <- s.Serve(lis)
		}()

		select {
		case err = <-errChan:
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				return fmt.Errorf("failed to serve: %w", err)
			}
			return nil
		case <-ctx.Done():
			slog.Info(fmt.Sprintf("stopping %s server", s.Name()), slog.String("addr", addr))
			stopCtx, stop := context.WithTimeout(context.Background(), shutdownTimeout)
			defer stop()
			return s.Stop(stopCtx)
		}
	})

	return nil
}

type ServerRunner interface {
	Run(ctx context.Context, g *errgroup.Group) error
}

func RunServers(ctx context.Context, runner ...ServerRunner) error {
	if len(runner) == 0 {
		return errors.New("no server runners")
	}

	g, ctx := errgroup.WithContext(ctx)

	for _, r := range runner {
		err := r.Run(ctx, g)
		if err != nil {
			return fmt.Errorf("failed to run server: %w", err)
		}
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("servers stopped with error: %w", err)
	}
	return nil
}

type GRPCRegistrator interface {
	RegisterServer(s *grpc.Server)
}

type grpcRunner struct {
	cfg       GRPCServerConfig
	registers []GRPCRegistrator
}

func NewGRPCRunner(cfg GRPCServerConfig, register ...GRPCRegistrator) ServerRunner {
	return &grpcRunner{
		cfg:       cfg,
		registers: append([]GRPCRegistrator{}, register...),
	}
}

func (r *grpcRunner) Run(ctx context.Context, g *errgroup.Group) error {
	grpcSrv := grpc.NewServer()

	for _, reg := range r.registers {
		reg.RegisterServer(grpcSrv)
	}

	reflection.Register(grpcSrv)

	return run(ctx, g, r.cfg.Host, r.cfg.Port, newGRPCAdapter("GRPC", grpcSrv))
}

type restRunner struct {
	cfg     RESTServerConfig
	handler http.Handler
}

func NewRESTServer(cfg RESTServerConfig, handler http.Handler) ServerRunner {
	return &restRunner{
		cfg:     cfg,
		handler: handler,
	}
}

func (r *restRunner) Run(ctx context.Context, g *errgroup.Group) error {
	restSrv := &http.Server{
		Handler:           r.handler,
		ReadHeaderTimeout: r.cfg.ReadHeaderTimeout,
	}

	return run(ctx, g, r.cfg.Host, r.cfg.Port, newHTTPAdapter("REST", restSrv))
}

type pprofRunner struct {
	cfg PprofConfig
}

func NewPprofRunner(cfg PprofConfig) ServerRunner {
	return &pprofRunner{cfg: cfg}
}

func (r *pprofRunner) Run(ctx context.Context, g *errgroup.Group) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	pprofSrv := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: r.cfg.ReadHeaderTimeout,
	}

	return run(ctx, g, r.cfg.Host, r.cfg.Port, newHTTPAdapter("Pprof", pprofSrv))
}
