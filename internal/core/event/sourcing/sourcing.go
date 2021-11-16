package sourcing

import (
	"context"
	"database/sql"
	"fmt"
	sql2 "github.com/cacos-group/cacos/internal/core/event/sourcing/sql"
	"github.com/cacos-group/cacos/internal/core/event/sourcing/strategy"
	"github.com/cacos-group/cacos/internal/core/leaf"
	"github.com/cacos-group/cacos/internal/core/metadata"
	clientV3 "go.etcd.io/etcd/client/v3"
)

const (
	AddNamespace = iota
	AddAppid
	AddKV
)

type Client interface {
	AddNamespace(ctx context.Context, namespace string) error
	AddAppid(ctx context.Context, namespace string, appid string) error
	AddKV(ctx context.Context, namespace string, appid string, name string, val string) error
}

type client struct {
	etcd     *clientV3.Client
	db       *sql.DB
	strategy strategy.Strategy

	leaf leaf.Leaf
}

func NewClient(db *sql.DB, etcd *clientV3.Client) Client {
	c := &client{
		etcd:     etcd,
		db:       db,
		strategy: strategy.New(db, etcd),
	}

	return c
}

func (c *client) AddNamespace(ctx context.Context, namespace string) (err error) {
	err = c.prepareEventLogStore(ctx, namespace)
	if err != nil {
		return
	}

	mds := metadata.Metadatas{}
	mds.Set(metadata.Namespace, namespace)

	err = c.strategy.Handler(ctx, AddNamespace, mds)
	if err != nil {
		return
	}

	return
}

func (c *client) AddAppid(ctx context.Context, namespace string, appid string) (err error) {
	mds := metadata.Metadatas{}
	mds.Set(metadata.Namespace, namespace)
	mds.Set(metadata.Appid, appid)

	err = c.strategy.Handler(ctx, AddAppid, mds)
	if err != nil {
		return
	}

	return
}

func (c *client) AddKV(ctx context.Context, namespace string, appid string, name string, val string) (err error) {
	key := fmt.Sprintf("/%s/%s/%s", namespace, appid, name)

	mds := metadata.Metadatas{}
	mds.Set(metadata.Namespace, namespace)
	mds.Set(metadata.Key, key)
	mds.Set(metadata.Val, val)

	err = c.strategy.Handler(ctx, AddKV, mds)
	if err != nil {
		return
	}

	return
}

func (c *client) prepareEventLogStore(ctx context.Context, namespace string) (err error) {
	conn, err := c.db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	tableName := strategy.GenTableName(namespace, 0)

	// create event_log_%s_%s
	_, err = conn.ExecContext(ctx, fmt.Sprintf(sql2.CreateEventLogSql, tableName))
	if err != nil {
		return
	}
	return
}
