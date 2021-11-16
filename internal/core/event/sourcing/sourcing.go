package sourcing

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/cacos-group/cacos/internal/core/event/sourcing/model"
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

type EventSourcing interface {
	GeneratorEvents(ctx context.Context, mds metadata.Metadatas) []model.Event
	Presentation(ctx context.Context, tableName string, events []model.Event) error
	Published(ctx context.Context, events []model.Event) (offset int, err error)
	Replayed(ctx context.Context, tableName string, events []model.Event, offset int) (isRetrySuccess bool, err error)
}

type client struct {
	etcd             *clientV3.Client
	db               *sql.DB
	eventSourcingMap map[int]EventSourcing

	leaf leaf.Leaf
}

func NewClient(db *sql.DB, etcd *clientV3.Client) Client {
	c := &client{
		etcd:             etcd,
		db:               db,
		eventSourcingMap: make(map[int]EventSourcing),
	}

	newStrategy := strategy.NewStrategy(db, etcd)
	c.eventSourcingMap[AddNamespace] = strategy.NewNamespace(newStrategy, db)
	c.eventSourcingMap[AddAppid] = strategy.NewAppid(newStrategy, db)
	c.eventSourcingMap[AddKV] = strategy.NewKV(newStrategy, db)

	return c
}

func (c *client) GetEventSourcing(name int) (EventSourcing, error) {
	es, ok := c.eventSourcingMap[name]
	if !ok {
		return nil, errors.New("event sourcing undefined")
	}

	return es.(EventSourcing), nil
}

func (c *client) AddNamespace(ctx context.Context, namespace string) (err error) {
	es, err := c.GetEventSourcing(AddNamespace)
	if err != nil {
		return
	}

	err = c.prepareEventLogStore(ctx, namespace)
	if err != nil {
		return
	}

	mds := metadata.Metadatas{}
	mds.Set(metadata.Namespace, namespace)

	events := es.GeneratorEvents(ctx, mds)

	tableName := strategy.GenTableName(namespace, 0)
	err = es.Presentation(ctx, tableName, events)
	if err != nil {
		return
	}

	err = c.published(es, ctx, tableName, events)
	if err != nil {
		return
	}
	return
}

func (c *client) AddAppid(ctx context.Context, namespace string, appid string) (err error) {
	es, err := c.GetEventSourcing(AddAppid)
	if err != nil {
		return
	}

	mds := metadata.Metadatas{}
	mds.Set(metadata.Namespace, namespace)
	mds.Set(metadata.Appid, appid)
	events := es.GeneratorEvents(ctx, mds)

	tableName := strategy.GenTableName(namespace, 0)
	err = es.Presentation(ctx, tableName, events)
	if err != nil {
		return
	}

	err = c.published(es, ctx, tableName, events)
	if err != nil {
		return
	}

	return nil
}

func (c *client) published(es EventSourcing, ctx context.Context, tableName string, events []model.Event) error {
	offset, err := es.Published(ctx, events)
	if err != nil {
		isRetrySuccess, replayedErr := es.Replayed(ctx, tableName, events, offset)
		if replayedErr != nil {
			return replayedErr
		}

		if isRetrySuccess == true {
			return nil
		}

		return err
	}

	return nil
}

func (c *client) AddKV(ctx context.Context, namespace string, appid string, name string, val string) (err error) {
	es, err := c.GetEventSourcing(AddKV)
	if err != nil {
		return
	}

	key := fmt.Sprintf("/%s/%s/%s", namespace, appid, name)

	mds := metadata.Metadatas{}
	mds.Set(metadata.Key, key)
	mds.Set(metadata.Val, val)
	events := es.GeneratorEvents(ctx, mds)

	tableName := strategy.GenTableName(namespace, 0)
	err = es.Presentation(ctx, tableName, events)
	if err != nil {
		return
	}

	err = c.published(es, ctx, tableName, events)
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
