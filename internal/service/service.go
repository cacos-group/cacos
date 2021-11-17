package service

import (
	"context"
	api "github.com/cacos-group/cacos/api/gen/go"
	"github.com/cacos-group/cacos/internal/core/event/sourcing"
	"github.com/cacos-group/cacos/internal/core/query"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/wire"
	"time"
)

var Provider = wire.NewSet(New, wire.Bind(new(api.CacosServer), new(*Service)))

//type Service cacosV1.CacosServer

// Service service.
type Service struct {
	api.UnimplementedCacosServer
	query         query.Client
	eventSourcing sourcing.Client
}

// New new a service and return.
func New(query query.Client, eventSourcing sourcing.Client) (s *Service, cf func(), err error) {
	return newService(query, eventSourcing)
}

func newService(query query.Client, eventSourcing sourcing.Client) (s *Service, cf func(), err error) {
	return &Service{
			query:         query,
			eventSourcing: eventSourcing,
		}, func() {

		}, nil
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

func (s *Service) NamespaceList(ctx context.Context, in *empty.Empty) (out *api.NamespaceListReply, err error) {
	out = new(api.NamespaceListReply)

	list, err := s.query.GetNamespaceList(ctx)
	if err != nil {
		return
	}

	out.NamespaceList = make([]*api.Namespace, 0, len(list))
	for _, item := range list {
		out.NamespaceList = append(out.NamespaceList, &api.Namespace{
			Namespace: item.Namespace,
		})
	}

	return out, nil
}

func (s *Service) AppList(ctx context.Context, in *api.AppListReq) (out *api.AppListReply, err error) {
	out = new(api.AppListReply)

	appList, err := s.query.GetAppidList(ctx, in.Namespace)
	if err != nil {
		return
	}

	out.AppList = make([]*api.App, 0, len(appList))
	for _, item := range appList {
		out.AppList = append(out.AppList, &api.App{
			App: item.Appid,
		})
	}

	return out, nil
}

func (s *Service) KvList(ctx context.Context, in *api.KVListReq) (out *api.KVListReply, err error) {
	out = new(api.KVListReply)

	list, err := s.query.GetKVList(ctx, in.Namespace, in.App)
	if err != nil {
		return
	}

	out.KvList = make([]*api.KV, 0, len(list.Kvs))
	for _, item := range list.Kvs {
		out.KvList = append(out.KvList, &api.KV{
			Namespace: in.Namespace,
			App:       in.App,
			Key:       string(item.Key),
			Val:       string(item.Value),
		})
	}

	return out, nil
}

func (s *Service) AddNamespace(ctx context.Context, in *api.AddNamespaceReq) (out *empty.Empty, err error) {
	err = s.eventSourcing.AddNamespace(ctx, in.Namespace)
	if err != nil {
		return
	}

	return nil, err
}

func (s *Service) AddApp(ctx context.Context, in *api.AddAppReq) (out *empty.Empty, err error) {
	err = s.eventSourcing.AddAppid(ctx, in.Namespace, in.App)
	if err != nil {
		return
	}

	return &empty.Empty{}, nil
}

func (s *Service) AddKV(ctx context.Context, in *api.AddKVReq) (out *empty.Empty, err error) {
	err = s.eventSourcing.AddKV(ctx, in.Namespace, in.App, in.Key, in.Val)
	if err != nil {
		return
	}

	return &empty.Empty{}, nil
}
