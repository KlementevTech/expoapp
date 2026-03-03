package servers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/sync/errgroup"
)

const shutdownTimeout = 5 * time.Second

type serverLauncher interface {
	Name() string
	Serve(listener net.Listener) error
	Stop(ctx context.Context) error
}

func launch(
	ctx context.Context,
	g *errgroup.Group,
	host string,
	port int,
	s serverLauncher,
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

func Run(ctx context.Context, runner ...Runner) error {
	if len(runner) == 0 {
		return errors.New("no server runners")
	}

	g, ctx := errgroup.WithContext(ctx)

	for _, r := range runner {
		err := r.Run(ctx, g)
		if err != nil {
			return fmt.Errorf("failed to launch server: %w", err)
		}
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("servers stopped with error: %w", err)
	}
	return nil
}
