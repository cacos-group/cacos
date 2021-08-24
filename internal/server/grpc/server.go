package grpc

import (
	"fmt"
	"github.com/cacos-group/cacos/internal/conf"
	"github.com/cacos-group/cacos/pkg/transport/grpc"
	"github.com/cacos-group/cacos/pkg/zaplog"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
)

func New(config *conf.Config, log zaplog.Logger) (s *grpc.Server, cf func(), err error) {
	return newServer(config, log)
}

func newServer(config *conf.Config, log zaplog.Logger) (s *grpc.Server, cf func(), err error) {
	s = grpc.NewServer(
		grpc.WithAddress(fmt.Sprintf(":%d", config.Server.Port)),
		grpc.WithLogger(log),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
			grpc_ctxtags.UnaryServerInterceptor())))

	cf = func() {

	}
	return
}
