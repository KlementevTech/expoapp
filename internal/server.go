package internal

import (
	"context"
	"errors"
	"fmt"
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

type servers struct {
	g     *errgroup.Group
	rs    *http.Server
	gs    *grpc.Server
	pprof *http.Server
}

func (s *servers) startRESTServer(cfg RESTServerConfig, handler http.Handler) {
	s.rs = &http.Server{
		Addr:              net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)),
		Handler:           handler,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
	}

	s.g.Go(func() error {
		slog.Default().Info(fmt.Sprintf("starting REST server at %s", s.rs.Addr))
		err := s.rs.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}

		slog.Default().Info("REST server stopped")
		return nil
	})
}

func (s *servers) startPprofServer(cfg PprofConfig) {
	s.pprof = &http.Server{
		Addr:              net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)),
		Handler:           http.DefaultServeMux,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
	}

	s.g.Go(func() error {
		slog.Default().Info(fmt.Sprintf("starting pprof server at %s", s.pprof.Addr))
		err := s.pprof.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}

		slog.Default().Info("pprof server stopped")
		return nil
	})
}

func (s *servers) startGRPCServer(cfg GRPCServerConfig, register ...func(*grpc.Server)) {
	s.gs = grpc.NewServer()

	for _, r := range register {
		r(s.gs)
	}

	reflection.Register(s.gs)

	s.g.Go(func() error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var lc net.ListenConfig
		lis, err := lc.Listen(ctx, "tcp", net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)))
		if err != nil {
			return err
		}

		slog.Default().InfoContext(ctx, fmt.Sprintf("starting gRPC server at %s", lis.Addr().String()))
		if err = s.gs.Serve(lis); err != nil {
			return err
		}

		slog.Default().InfoContext(ctx, "gRPC server stopped")
		return nil
	})
}

const shutdownTimeout = 5 * time.Second

func (s *servers) shutdownServers(ctx context.Context) {
	s.g.Go(func() error {
		<-ctx.Done()
		go s.gs.GracefulStop()
		go func() {
			shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
			defer cancel()

			err := s.rs.Shutdown(shutdownCtx)
			if err != nil {
				slog.Default().Error(fmt.Sprintf("gRPC server stopped with error: %v", err))
			}
		}()

		if s.pprof != nil {
			go func() {
				shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
				defer cancel()

				err := s.pprof.Shutdown(shutdownCtx)
				if err != nil {
					slog.Default().Error(fmt.Sprintf("pprof server stopped with error: %v", err))
				}
			}()
		}

		return nil
	})
}

func (s *servers) wait() error {
	return s.g.Wait()
}

func RunServers(ctx context.Context, cfg Config, vs *service.VersionService) error {
	g, ctx := errgroup.WithContext(ctx)
	s := &servers{
		g: g,
	}

	if cfg.Pprof.Enabled {
		s.startPprofServer(cfg.Pprof)
	}
	s.startRESTServer(cfg.RESTServer, rest.SetupRouter(vs))
	s.startGRPCServer(cfg.GRPCServer, expo.Register(vs))
	s.shutdownServers(ctx)

	return s.wait()
}
