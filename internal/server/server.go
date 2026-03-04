package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"strconv"
	"time"

	"golang.org/x/sync/errgroup"
)

const shutdownTimeout = 5 * time.Second

type ServeStopAdapter interface {
	Serve(listener net.Listener) error
	Stop(ctx context.Context) error
}
type Server struct {
	name    string
	address string
	adapter ServeStopAdapter
}

func NewServer(name string, host string, port int, adapter ServeStopAdapter) Server {
	return Server{
		name:    name,
		address: net.JoinHostPort(host, strconv.Itoa(port)),
		adapter: adapter,
	}
}

func (s *Server) Run(
	ctx context.Context,
	g *errgroup.Group,
) error {
	var lc net.ListenConfig
	lis, err := lc.Listen(ctx, "tcp4", s.address)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	slog.Default().InfoContext(ctx, fmt.Sprintf("%s server is listening", s.name), slog.String("addr", s.address))

	g.Go(func() error {
		errChan := make(chan error, 1)
		defer close(errChan)

		go func() {
			serveErr := s.adapter.Serve(lis)
			if serveErr != nil {
				errChan <- serveErr
			}
		}()

		select {
		case err = <-errChan:
			return fmt.Errorf("failed to serve: %w", err)
		case <-ctx.Done():
			slog.Info(fmt.Sprintf("stopping %s server", s.name), slog.String("addr", s.address))
			stopCtx, stop := context.WithTimeout(context.Background(), shutdownTimeout)
			defer stop()
			return s.adapter.Stop(stopCtx)
		}
	})

	return nil
}

type Runner interface {
	Run(ctx context.Context, g *errgroup.Group) error
}

func Run(ctx context.Context, runner ...Runner) error {
	if len(runner) == 0 {
		return errors.New("no server runners")
	}

	g, gCtx := errgroup.WithContext(ctx)

	for _, r := range runner {
		err := r.Run(gCtx, g)
		if err != nil {
			return fmt.Errorf("failed to launch server: %w", err)
		}
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("servers stopped with error: %w", err)
	}
	return nil
}
