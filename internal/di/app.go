package di

import (
	"fmt"
	api "github.com/cacos-group/cacos/api"
	"github.com/cacos-group/cacos/internal/conf"
	"github.com/cacos-group/cacos/internal/service"
	"google.golang.org/grpc"
	"log"
	"net"
)

type App struct {
	config *conf.Config
	svc    *service.Service
	grpc   *grpc.Server
}

func NewApp(config *conf.Config, svc *service.Service, g *grpc.Server) (app *App, closeFunc func(), err error) {
	app = &App{
		config: config,
		svc:    svc,
		grpc:   g,
	}
	closeFunc = func() {
		g.Stop()
		log.Printf("grpc.Stop done")
	}
	return
}

func (app *App) Start() error {
	serverConfig := app.config.Server
	port := serverConfig.Port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	api.RegisterCacosServer(app.grpc, app.svc)
	err = app.grpc.Serve(lis)
	if err != nil {
		return err
	}

	return nil
}
