package sourcing

import (
	"context"
	"database/sql"
	"github.com/cacos-group/cacos/internal/core/event/sourcing/controller"
	"github.com/cacos-group/cacos/internal/core/event/sourcing/model"
	"github.com/cacos-group/cacos/internal/core/leaf"
	clientV3 "go.etcd.io/etcd/client/v3"
)

type Client interface {
	AddNamespace(ctx context.Context, namespace string) error
	AddAppid(ctx context.Context, namespace string, appid string) error
	AddKV(ctx context.Context, namespace string, appid string, name string, val string) error
}

type client struct {
	etcd *clientV3.Client
	db   *sql.DB

	leaf leaf.Leaf

	controller controller.Controller
}

func NewClient(db *sql.DB, etcd *clientV3.Client) Client {
	c := &client{
		etcd:       etcd,
		db:         db,
		controller: controller.New(db, etcd),
	}

	return c
}

func (c *client) AddNamespace(ctx context.Context, namespace string) (err error) {
	err = c.controller.Handler(ctx, model.GeneratorNamespaceReq(namespace))
	if err != nil {
		return
	}

	return
}

func (c *client) AddAppid(ctx context.Context, namespace string, appid string) (err error) {
	err = c.controller.Handler(ctx, model.GeneratorServiceReq(namespace, appid))
	if err != nil {
		return
	}

	return
}

func (c *client) AddKV(ctx context.Context, namespace string, appid string, name string, val string) (err error) {
	err = c.controller.Handler(ctx, model.GeneratorKvReq(namespace, appid, name, val))
	if err != nil {
		return
	}

	return
}
