package internal

type Config struct {
	GRPCServerHost string `env:"GRPC_SERVER_HOST" envDefault:"127.0.0.1"`
	GRPCServerPort int    `env:"GRPC_SERVER_PORT" envDefault:"50051"`
}
