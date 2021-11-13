package sourcing

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/cacos-group/cacos/internal/core/event/sourcing/model"
	"github.com/cacos-group/cacos/internal/core/event/sourcing/strategy"
	"github.com/cacos-group/cacos/internal/core/leaf"
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
	Prepare(ctx context.Context, namespace string, appid string) error
	Presentation(ctx context.Context, tableName string, events []model.Event) error
	Published(ctx context.Context, events []model.Event) error
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

	err = es.Prepare(ctx, namespace, "")
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

	err = es.Prepare(ctx, namespace, appid)
	if err != nil {
		return
	}

	key := fmt.Sprintf("/%s/%s", namespace, appid)

	events := make([]model.Event, 0, 5)

	// event KVPut
	events = append(events, model.NewAppidPutEvent(key, ""))

	user := fmt.Sprintf("u_%s_%s", namespace, appid)
	// todo 生成password 和 加密
	password := "password"

	// event UserAdd
	events = append(events, model.NewUserAddEvent(key, user, password))

	// event RoleAdd
	role := fmt.Sprintf("%s_%s", namespace, appid)
	events = append(events, model.NewRoleAddEvent(key, role))

	// event UserGrantRole
	events = append(events, model.NewUserGrantRoleEvent(key, user, role))

	// event RoleGrantPermission
	//todo 权限控制读写分离
	permissionType := "2"
	events = append(events, model.NewRoleGrantPermissionEvent(role, key, permissionType))

	tableName := strategy.GenTableName(namespace, appid)
	err = es.Presentation(ctx, tableName, events)
	if err != nil {
		return
	}

	err = es.Published(ctx, events)
	if err != nil {
		return
	}

	return nil
}

func (c *client) AddKV(ctx context.Context, namespace string, appid string, name string, val string) (err error) {
	es, err := c.GetEventSourcing(AddKV)
	if err != nil {
		return
	}

	err = es.Prepare(ctx, namespace, appid)
	if err != nil {
		return
	}

	key := fmt.Sprintf("/%s/%s/%s", namespace, appid, name)

	events := make([]model.Event, 0)
	// event KVPut
	events = append(events, model.NewKVPutEvent(key, val))

	tableName := strategy.GenTableName(namespace, appid)
	err = es.Presentation(ctx, tableName, events)
	if err != nil {
		return
	}

	err = es.Published(ctx, events)
	if err != nil {
		return
	}

	return
}
