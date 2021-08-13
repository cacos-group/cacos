package resource

import (
	"github.com/cacos-group/cacos-server-sdk/api"
	"github.com/cacos-group/cacos-server-sdk/entry"
	"github.com/cacos-group/cacos/internal/conf"
	"time"
)

func NewCacos(config *conf.Config) (c entry.Cacos, cf func(), err error) {
	mysqlConfig := config.Mysql
	etcdConfig := config.Etcd
	return entry.New(&api.Config{
		Mysql: api.MysqlConfig{
			DSN:             mysqlConfig.DSN,
			MaxOpenConns:    mysqlConfig.MaxOpenConns,
			MaxIdleConns:    mysqlConfig.MaxIdleConns,
			ConnMaxLifetime: time.Duration(mysqlConfig.ConnMaxLifetime),
			ConnMaxIdleTime: time.Duration(mysqlConfig.ConnMaxIdleTime),
		},
		Etcd: etcdConfig,
	})
}
