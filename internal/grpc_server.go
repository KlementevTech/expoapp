package internal

import (
	"context"
	"net"
	"strconv"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type RegisterFunc func(s *grpc.Server)

func StartGRPCServer(
	ctx context.Context,
	g *errgroup.Group,
	host string,
	port int,
	register ...RegisterFunc,
) *grpc.Server {
	s := grpc.NewServer()
	reflection.Register(s)

	for _, r := range register {
		r(s)
	}

	g.Go(func() error {
		lis, err := new(net.ListenConfig).
			Listen(ctx, "tcp", net.JoinHostPort(host, strconv.Itoa(port)))
		if err != nil {
			return err
		}

		defer func() {
			_ = lis.Close()
		}()

		if err = s.Serve(lis); err != nil {
			return err
		}

		return nil
	})

	return s
}
