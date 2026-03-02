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

type Runner interface {
	Run(ctx context.Context, g *errgroup.Group) error
}

func RunServers(ctx context.Context, runner ...Runner) error {
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

type RegisterFunc func(s *grpc.Server)

type GRPCRunner struct {
	cfg       GRPCServerConfig
	registers []RegisterFunc
}

func NewGRPCRunner(cfg GRPCServerConfig, register ...RegisterFunc) *GRPCRunner {
	return &GRPCRunner{
		cfg:       cfg,
		registers: append([]RegisterFunc{}, register...),
	}
}

func (r *GRPCRunner) Run(ctx context.Context, g *errgroup.Group) error {
	grpcSrv := grpc.NewServer()

	for _, register := range r.registers {
		register(grpcSrv)
	}

	reflection.Register(grpcSrv)

	return run(ctx, g, r.cfg.Host, r.cfg.Port, newGRPCAdapter("GRPC", grpcSrv))
}

type RESTRunner struct {
	cfg     RESTServerConfig
	handler http.Handler
}

func NewRESTServer(cfg RESTServerConfig, handler http.Handler) *RESTRunner {
	return &RESTRunner{
		cfg:     cfg,
		handler: handler,
	}
}

func (r *RESTRunner) Run(ctx context.Context, g *errgroup.Group) error {
	restSrv := &http.Server{
		Handler:           r.handler,
		ReadHeaderTimeout: r.cfg.ReadHeaderTimeout,
	}

	return run(ctx, g, r.cfg.Host, r.cfg.Port, newHTTPAdapter("REST", restSrv))
}

type PprofRunner struct {
	cfg PprofConfig
}

func NewPprofRunner(cfg PprofConfig) *PprofRunner {
	return &PprofRunner{cfg: cfg}
}

func (r *PprofRunner) Run(ctx context.Context, g *errgroup.Group) error {
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
