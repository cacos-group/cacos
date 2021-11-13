package resource

import (
	"context"
	"database/sql"
	"github.com/cacos-group/cacos/internal/conf"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"log"
)

func NewDB(cfg *conf.Config) (db *sql.DB, cf func(), err error) {
	log.Println("NewDB start")

	mysqlConfig := cfg.Mysql
	db, err = sql.Open("mysql", mysqlConfig.DSN)
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(mysqlConfig.ConnMaxIdleTime))
	if err = db.PingContext(ctx); err != nil {
		return
	}
	db.SetMaxOpenConns(mysqlConfig.MaxOpenConns)
	db.SetMaxIdleConns(mysqlConfig.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(mysqlConfig.ConnMaxLifetime))
	db.SetConnMaxIdleTime(time.Duration(mysqlConfig.ConnMaxIdleTime))

	cf = func() {
		db.Close()
	}

	log.Println("NewDB done")

	return
}
