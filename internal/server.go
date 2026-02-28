package internal

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	_ "net/http/pprof" //nolint:gosec // pprof нужен для профилирования в dev-окружении
	"strconv"
	"time"

	"expo/internal/handlers/expo"
	"expo/internal/rest"
	"expo/internal/service"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const shutdownTimeout = 5 * time.Second

type stopFunc func()

func startRESTServer(g *errgroup.Group, cfg RESTServerConfig, handler http.Handler) stopFunc {
	s := &http.Server{
		Addr:              net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)),
		Handler:           handler,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
	}

	stop := func() {
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		err := s.Shutdown(ctx)
		if err != nil {
			slog.Default().Error("failed on shutdown REST server", "error", err)
		}
	}

	g.Go(func() error {
		slog.Default().Info("starting REST server", slog.String("addr", s.Addr))
		err := s.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	return stop
}

func startPprofServer(g *errgroup.Group, cfg PprofConfig) stopFunc {
	s := &http.Server{
		Addr:              net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)),
		Handler:           http.DefaultServeMux,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
	}

	stop := func() {
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		err := s.Shutdown(ctx)
		if err != nil {
			slog.Default().Error("failed on shutdown Pprof server", "error", err)
		}
	}

	g.Go(func() error {
		slog.Default().Info("starting Pprof server", slog.String("addr", s.Addr))
		err := s.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	return stop
}

func startGRPCServer(g *errgroup.Group, cfg GRPCServerConfig, register ...func(*grpc.Server)) stopFunc {
	s := grpc.NewServer()
	for _, r := range register {
		r(s)
	}

	reflection.Register(s)

	stop := func() {
		s.GracefulStop()
	}

	g.Go(func() error {
		var lc net.ListenConfig
		lis, err := lc.Listen(context.Background(), "tcp4", net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)))
		if err != nil {
			return err
		}
		defer func() {
			_ = lis.Close()
		}()

		slog.Default().Info("starting gRPC server", slog.String("addr", lis.Addr().String()))
		if err = s.Serve(lis); err != nil {
			return err
		}
		return nil
	})

	return stop
}

func RunServers(ctx context.Context, cfg *Config, vs *service.VersionService) error {
	g, ctx := errgroup.WithContext(ctx)

	stopGrpc := startGRPCServer(g, cfg.GRPCServer, expo.Register(vs))
	stopRest := startRESTServer(g, cfg.RESTServer, rest.SetupRouter(vs))

	var stopPprof stopFunc
	if cfg.Pprof.Enabled {
		stopPprof = startPprofServer(g, cfg.Pprof)
	}

	g.Go(func() error {
		<-ctx.Done()

		stopGrpc()
		stopRest()

		if stopPprof != nil {
			stopPprof()
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}

	slog.Default().InfoContext(ctx, "servers shutdown gracefully")
	return nil
}
