package web

import (
	"context"
	"net"
	"strconv"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type ServiceRegister interface {
	RegisterService(s *grpc.Server)
}

type ServiceRegisterFunc func(s *grpc.Server)

func (f ServiceRegisterFunc) RegisterService(s *grpc.Server) {
	f(s)
}

func NewGRPCServer(register ...ServiceRegister) *grpc.Server {
	s := grpc.NewServer()

	for _, r := range register {
		r.RegisterService(s)
	}

	reflection.Register(s)

	return s
}

func StartGRPCServer(
	ctx context.Context,
	g *errgroup.Group,
	s *grpc.Server,
	host string,
	port int,
) *grpc.Server {
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
