package server

import (
	v1 "nancalacc/api/account/v1"
	"nancalacc/internal/conf"
	"nancalacc/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"

	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, accountService *service.AccountService, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	healthCheck := health.NewServer()
	healthCheck.SetServingStatus("nancalacc.account.v1.asyncacc", grpc_health_v1.HealthCheckResponse_SERVING)
	v1.RegisterAccountServer(srv, accountService)
	return srv
}
