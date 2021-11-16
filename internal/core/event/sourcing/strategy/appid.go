package strategy

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/cacos-group/cacos/internal/core/event/sourcing/model"
	"github.com/cacos-group/cacos/internal/core/metadata"
	clientV3 "go.etcd.io/etcd/client/v3"
)

type Appid struct {
	*strategy
	db *sql.DB
}

func NewAppid(strategy *strategy, db *sql.DB) *Appid {
	n := &Appid{}
	n.strategy = strategy
	n.db = db

	return n
}

func (s *Appid) GeneratorEvents(ctx context.Context, mds metadata.Metadatas) (list []model.Event) {
	namespace := mds.Get(metadata.Namespace)
	appid := mds.Get(metadata.Appid)

	key := fmt.Sprintf("/%s/%s", namespace, appid)

	user := fmt.Sprintf("u_%s_%s", namespace, appid)
	// todo 生成password 和 加密
	password := "password"

	role := fmt.Sprintf("%s_%s", namespace, appid)

	// 只读权限
	permissionType := string(clientV3.PermRead)

	list = []model.Event{
		model.NewInfoAppidPutEvent(namespace, appid),                 //
		model.NewAppidPutEvent(key, ""),                              // event KVPut
		model.NewUserAddEvent(key, user, password),                   // event UserAdd
		model.NewRoleAddEvent(key, role),                             // event RoleAdd
		model.NewUserGrantRoleEvent(key, user, role),                 // event UserGrantRole
		model.NewRoleGrantPermissionEvent(role, key, permissionType), // // event RoleGrantPermission
	}

	return
}
