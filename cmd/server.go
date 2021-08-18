package main

import (
	"context"
	"flag"
	"github.com/cacos-group/cacos/internal/conf"
	"github.com/cacos-group/cacos/internal/di"
	"github.com/cacos-group/cacos/pkg/transport/http"
	"github.com/cacos-group/cacos/pkg/zaplog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	flag.Parse()

	config := conf.Init()

	app, cf, err := di.InitApp(config)
	if err != nil {
		panic(err)
	}

	http.NewServer().Start(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		app.Log().Info("signal.Notify", zaplog.Any("signal", s.String()))
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			cf()
			app.Log().Info("exit", zaplog.Any("signal", s.String()))
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
