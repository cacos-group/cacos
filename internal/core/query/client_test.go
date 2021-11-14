package query

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

	cfg.Mysql = conf.MysqlConfig{
		DSN:             "admin:admin@tcp(127.0.0.1:3306)/cacos",
		ConnMaxLifetime: conf.Duration(60 * time.Second),
		ConnMaxIdleTime: conf.Duration(6 * time.Hour),
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

func TestClient_GetNamespaceList(t *testing.T) {
	c := testClient(t)
	list, err := c.GetNamespaceList(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(list)
}

func TestClient_GetAppidList(t *testing.T) {
	c := testClient(t)
	list, err := c.GetAppidList(context.Background(), "namespace5")
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(list)
}

func TestClient_GetKVList(t *testing.T) {
	c := testClient(t)
	list, err := c.GetKVList(context.Background(), "namespace5", "appid9")
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(list)
}
