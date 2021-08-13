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

func (s *Service) AuthLogin(ctx context.Context, req *api.LoginRequest) (reply *api.LoginReply, err error) {
	return &api.LoginReply{
		Token: "asdasd",
	}, nil
}

func (s *Service) NamespaceList(ctx context.Context, in *api.NamespaceListReq) (out *api.NamespaceListReply, err error) {
	out = new(api.NamespaceListReply)

	a, err := s.cacos.Administer(ctx)
	if err != nil {
		return
	}

	namespaceList, err := a.GetNamespaceList(ctx)
	if err != nil {
		return
	}

	out.NamespaceList = make([]*api.Namespace, 0, len(namespaceList))
	for _, item := range namespaceList {
		out.NamespaceList = append(out.NamespaceList, &api.Namespace{
			Namespace: item.Namespace,
		})
	}

	return out, nil
}

func (s *Service) AppList(ctx context.Context, in *api.AppListReq) (out *api.AppListReply, err error) {
	out = new(api.AppListReply)

	a, err := s.cacos.Administer(ctx)
	if err != nil {
		return
	}

	appList, err := a.GetAppList(ctx, in.Namespace)
	if err != nil {
		return
	}

	out.AppList = make([]*api.App, 0, len(appList))
	for _, item := range appList {
		out.AppList = append(out.AppList, &api.App{
			Namespace: item.Namespace,
			App:       item.App,
		})
	}

	return out, nil
}
