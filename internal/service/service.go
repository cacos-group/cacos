package service

import (
	"context"
	"github.com/cacos-group/cacos-server-sdk/entry"
	api "github.com/cacos-group/cacos/api"
	"github.com/google/wire"
	"time"
)

var Provider = wire.NewSet(New, wire.Bind(new(api.CacosServer), new(*Service)))

//type Service cacosV1.CacosServer

// Service service.
type Service struct {
	api.UnimplementedCacosServer
	cacos entry.Cacos
}

// New new a service and return.
func New(cacos entry.Cacos) (s *Service, cf func(), err error) {
	return newService(cacos)
}

func newService(cacos entry.Cacos) (s *Service, cf func(), err error) {
	return &Service{
		cacos: cacos,
	}, nil, nil
}

func (s *Service) SayHello(ctx context.Context, req *api.HelloRequest) (reply *api.HelloReply, err error) {
	time.Sleep(4 * time.Second)
	return &api.HelloReply{
		Message: req.Name,
	}, nil
}
