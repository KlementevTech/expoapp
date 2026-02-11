package web

import (
	"context"
	"log/slog"
	"net"

	expopb "expoapp/pkg/api/expo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func StartGRPCServer(
	ctx context.Context,
	addr string,
	expoService expopb.ExpoServiceServer,
) (*grpc.Server, <-chan error) {
	s := grpc.NewServer()

	expopb.RegisterExpoServiceServer(s, expoService)

	reflection.Register(s)

	errChan := make(chan error, 1)

	go func() {
		defer close(errChan)

		lis, err := (&net.ListenConfig{}).Listen(ctx, "tcp", addr)
		if err != nil {
			errChan <- err
			return
		}
		defer func() {
			if lisErr := lis.Close(); lisErr != nil {
				slog.Default().Error("failed to close listener", "error", lisErr)
			}
		}()

		if err = s.Serve(lis); err != nil {
			errChan <- err
			return
		}
	}()

	return s, errChan
}
