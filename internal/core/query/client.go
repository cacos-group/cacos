package query

import (
	"context"
	"database/sql"
	"fmt"
	model2 "github.com/cacos-group/cacos/internal/core/event/sourcing/model"
	"github.com/cacos-group/cacos/internal/core/query/model"
	clientV3 "go.etcd.io/etcd/client/v3"
)

type Client interface {
	GetNamespaceList(ctx context.Context) (list []model.NamespaceModel, err error)
	GetAppidList(ctx context.Context, namespace string) (list []model.AppidModel, err error)
	GetKVList(ctx context.Context, namespace string, appid string) (rsp *clientV3.GetResponse, err error)
}

type client struct {
	db   *sql.DB
	etcd *clientV3.Client
}

func NewClient(db *sql.DB, etcd *clientV3.Client) Client {
	c := &client{
		etcd: etcd,
		db:   db,
	}

	return c
}

func (c *client) GetNamespaceList(ctx context.Context) (list []model.NamespaceModel, err error) {
	prefix := model2.GenNamespacePrefix()
	var opts []clientV3.OpOption
	opts = append(opts, clientV3.WithPrefix())

	rsp, err := c.etcd.Get(ctx, prefix, opts...)
	if err != nil {
		return
	}

	list = make([]model.NamespaceModel, 0, rsp.Count)
	for _, item := range rsp.Kvs {
		list = append(list, model.NamespaceModel{Namespace: string(item.Value)})
	}

	return
}

func (c *client) GetAppidList(ctx context.Context, namespace string) (list []model.AppidModel, err error) {
	prefix := model2.GenAppidPrefix(namespace)
	var opts []clientV3.OpOption
	opts = append(opts, clientV3.WithPrefix())

	rsp, err := c.etcd.Get(ctx, prefix, opts...)
	if err != nil {
		return
	}

	list = make([]model.AppidModel, 0, rsp.Count)
	for _, item := range rsp.Kvs {
		list = append(list, model.AppidModel{Appid: string(item.Value)})
	}

	return
}

func (c *client) GetKVList(ctx context.Context, namespace string, appid string) (rsp *clientV3.GetResponse, err error) {
	path := fmt.Sprintf("/%s/%s/", namespace, appid)
	var opts []clientV3.OpOption
	opts = append(opts, clientV3.WithPrefix())

	rsp, err = c.etcd.Get(ctx, path, opts...)
	if err != nil {
		return
	}

	//list = make([]model.KVModel, 0, 20)
	//for _, item := range rsp.Kvs {
	//	newKV := model.KVModel{
	//		Key:            string(item.Key),
	//		CreateRevision: item.CreateRevision,
	//		ModRevision:    item.ModRevision,
	//		Version:        item.Version,
	//		Value:          string(item.Value),
	//	}
	//	list = append(list, newKV)
	//}

	return
}
