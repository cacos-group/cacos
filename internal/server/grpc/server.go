package grpc

import (
	"fmt"
	"github.com/cacos-group/cacos/api"
	"github.com/cacos-group/cacos/internal/conf"
	"github.com/cacos-group/cacos/internal/service"
	"github.com/cacos-group/cacos/pkg/zaplog"
	"google.golang.org/grpc"
	"net"
)

type Server interface {
	Stop()
}

type server struct {
	Server *grpc.Server
}

func New(config *conf.Config, svc *service.Service, g *grpc.Server, log zaplog.Logger) (s Server, cf func(), err error) {
	return newServer(config, svc, g, log)
}

func newServer(config *conf.Config, svc *service.Service, g *grpc.Server, log zaplog.Logger) (s Server, cf func(), err error) {
	serverConfig := config.Server
	port := serverConfig.Port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(fmt.Sprintf("failed to listen: %v", err))
	}

	apiV1.RegisterCacosServer(g, svc)
	go func() {
		err = g.Serve(lis)
		if err != nil {
			panic(err)
		}
	}()

	cf = func() {
		g.Stop()
	}

	return &server{Server: g}, cf, nil
}

func (s *server) Stop() {
	s.Stop()
}
