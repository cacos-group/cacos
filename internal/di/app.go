package di

import (
	"fmt"
	api "github.com/cacos-group/cacos/api"
	"github.com/cacos-group/cacos/internal/conf"
	"github.com/cacos-group/cacos/internal/service"
	zaplog "github.com/cacos-group/cacos/pkg/zaplog"
	"google.golang.org/grpc"
	"net"
)

type App struct {
	config *conf.Config
	svc    *service.Service
	grpc   *grpc.Server
	log    zaplog.Logger
}

func NewApp(config *conf.Config, svc *service.Service, g *grpc.Server, log zaplog.Logger) (app *App, closeFunc func(), err error) {
	app = &App{
		config: config,
		svc:    svc,
		grpc:   g,
		log:    log,
	}
	closeFunc = func() {
		g.Stop()
		log.Info("app stop")
	}
	log.Info("app start")
	return
}

func (app *App) Start() error {
	serverConfig := app.config.Server
	port := serverConfig.Port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		app.log.Fatal(fmt.Sprintf("failed to listen: %v", err))
		return err
	}

	api.RegisterCacosServer(app.grpc, app.svc)
	err = app.grpc.Serve(lis)
	if err != nil {
		return err
	}

	return nil
}
