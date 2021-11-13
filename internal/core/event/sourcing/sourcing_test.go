package sourcing

import (
	"context"
	"github.com/cacos-group/cacos/internal/core/conf"
	"github.com/cacos-group/cacos/internal/core/resource"
	"testing"
	"time"
)

func testClient(t *testing.T) Client {
	cfg := new(conf.Config)

	cfg.Mysql = conf.MysqlConfig{
		DSN:             "admin:admin@tcp(127.0.0.1:3306)/cacos",
		ConnMaxLifetime: 60 * time.Second,
		ConnMaxIdleTime: 6 * time.Hour,
	}

	db, _, err := resource.NewDB(cfg)
	if err != nil {
		t.Error(err)
		return nil
	}

	cfg.Etcd = conf.EtcdConfig{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
		Username:    "root",
		Password:    "",
	}

	c, _, err := resource.NewEtcd(cfg)
	if err != nil {
		t.Error(err)
		return nil
	}

	s := NewClient(db, c)

	return s
}

func TestClient_AddNamespace(t *testing.T) {
	c := testClient(t)
	err := c.AddNamespace(context.Background(), "namespace5")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestClient_AddAppid(t *testing.T) {
	c := testClient(t)

	err := c.AddAppid(context.Background(), "namespace5", "appid9")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestClient_AddKV(t *testing.T) {
	c := testClient(t)

	err := c.AddKV(context.Background(), "namespace5", "appid9", "key1", "val2")
	if err != nil {
		t.Error(err)
		return
	}
}
