package resource

import (
	"github.com/cacos-group/cacos/internal/conf"
	clientV3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

func NewEtcd(cfg *conf.Config) (c *clientV3.Client, cf func(), err error) {
	log.Println("NewEtcd start")
	etcdCfg := cfg.Etcd

	c, err = clientV3.New(clientV3.Config{
		Endpoints:   etcdCfg.Endpoints,
		DialTimeout: 5 * time.Second,
		Username:    etcdCfg.Username,
		Password:    etcdCfg.Password,
	})
	if err != nil {
		return nil, cf, err
	}

	cf = func() {
		c.Close()
	}

	log.Println("NewEtcd done")
	return c, cf, nil
}
