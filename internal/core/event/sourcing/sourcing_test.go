package sourcing

import (
	"context"
	"fmt"
	"github.com/cacos-group/cacos/internal/conf"
	"github.com/cacos-group/cacos/internal/resource"
	"testing"
	"time"
)

func testClient(t *testing.T) Client {
	cfg := new(conf.Config)

	cfg.Mysql.DSN = "admin:admin@tcp(127.0.0.1:3306)/cacos"
	cfg.Mysql.ConnMaxLifetime = conf.Duration(60 * time.Second)
	cfg.Mysql.ConnMaxIdleTime = conf.Duration(6 * time.Hour)

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
	err := c.AddNamespace(context.Background(), fmt.Sprintf("namespace%d", time.Now().Unix()))
	if err != nil {
		t.Error(err)
		return
	}
}

func TestClient_AddAppid(t *testing.T) {
	c := testClient(t)

	err := c.AddNamespace(context.Background(), "namespace23")
	err = c.AddAppid(context.Background(), "namespace23", fmt.Sprintf("appid%d", time.Now().Unix()))
	if err != nil {
		t.Error(err)
		return
	}
}

func TestClient_AddKV(t *testing.T) {
	c := testClient(t)

	err := c.AddKV(context.Background(), "namespace5", "appid9", fmt.Sprintf("key%d", time.Now().Unix()), fmt.Sprintf("val%d", time.Now().Unix()))
	if err != nil {
		t.Error(err)
		return
	}
}
