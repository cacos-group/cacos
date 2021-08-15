package di

import (
	"github.com/cacos-group/cacos/internal/conf"
	grpc2 "github.com/cacos-group/cacos/internal/server/grpc"
	"github.com/cacos-group/cacos/internal/service"
	zaplog "github.com/cacos-group/cacos/pkg/zaplog"
)

type App struct {
	config *conf.Config
	svc    *service.Service
	grpc   grpc2.Server
	log    zaplog.Logger
}

func NewApp(config *conf.Config, svc *service.Service, g grpc2.Server, log zaplog.Logger) (app *App, closeFunc func(), err error) {
	app = &App{
		config: config,
		svc:    svc,
		grpc:   g,
		log:    log,
	}
	closeFunc = func() {
		log.Info("app stop")
	}
	log.Info("app start")
	return
}

func (app *App) Log() zaplog.Logger {
	return app.log
}
