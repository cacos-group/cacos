package resource

import (
	"github.com/cacos-group/cacos/internal/conf"
	"google.golang.org/grpc"
)

func NewGRPCServer(config *conf.Config) *grpc.Server {
	return grpc.NewServer()
}
