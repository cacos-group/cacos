package resource

import (
	"github.com/cacos-group/cacos/internal/conf"
	grpc_timeout "github.com/cacos-group/cacos/pkg/go-grpc-middleware/timeout"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc"
	"time"
)

func NewGRPCServer(cfg *conf.Config) *grpc.Server {
	return grpc.NewServer(

		grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(
				grpc_recovery.StreamServerInterceptor(),
				grpc_ctxtags.StreamServerInterceptor())),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_timeout.UnaryServerInterceptor(time.Duration(cfg.Server.Timeout)),
			grpc_recovery.UnaryServerInterceptor(),
			grpc_ctxtags.UnaryServerInterceptor())))
}
