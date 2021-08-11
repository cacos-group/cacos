package conf

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/cacos-group/cacos-server-sdk/api"
	"log"
	"time"
)

var (
	confPath string
)

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
}

type Config struct {
	Etcd   api.EtcdConfig
	Mysql  mysql
	Server server
}

type server struct {
	Port int
}

type etcd struct {
}

type mysql struct {
	DSN             string   // data source name.
	MaxOpenConns    int      // pool
	MaxIdleConns    int      // pool
	ConnMaxLifetime Duration //
	ConnMaxIdleTime Duration //
}

// Duration be used toml unmarshal string time, like 1s, 500ms.
type Duration time.Duration

// UnmarshalText unmarshal text to duration.
func (d *Duration) UnmarshalText(text []byte) error {
	tmp, err := time.ParseDuration(string(text))
	if err == nil {
		*d = Duration(tmp)
	}
	return err
}

func Init() (config *Config) {
	if _, err := toml.DecodeFile(confPath, &config); err != nil {
		log.Fatalf("err: %v", err)
		return
	}

	return
}