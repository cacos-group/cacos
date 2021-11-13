package di

import (
	"context"
	"fmt"
	apiV1 "github.com/cacos-group/cacos/api/gen/go"
	gw "github.com/cacos-group/cacos/api/gen/go"
	"github.com/cacos-group/cacos/internal/conf"
	"github.com/cacos-group/cacos/internal/service"
	"github.com/cacos-group/cacos/pkg/transport/grpc"
	zaplog "github.com/cacos-group/cacos/pkg/zaplog"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	xgrpc "google.golang.org/grpc"
	"net/http"
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

	apiV1.RegisterCacosServer(g, svc)
	go func() {
		err = g.Start(context.Background())
		if err != nil {
			panic(err)
		}
		log.Info("http server start")
	}()

	go func() {
		err = app.startGrpcGateway(config)
		if err != nil {
			panic(err)
		}
		log.Info("grpc gatewaty server start")
	}()

	closeFunc = func() {
		g.Stop(context.Background())
		log.Info("app stop")
	}
	log.Info("app start")
	return
}

func (app *App) Log() zaplog.Logger {
	return app.log
}

func (app *App) startGrpcGateway(config *conf.Config) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	opts := []xgrpc.DialOption{xgrpc.WithInsecure()}
	err := gw.RegisterCacosHandlerFromEndpoint(ctx, mux, fmt.Sprintf(":%d", config.Server.Port), opts)
	if err != nil {
		return err
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return http.ListenAndServe(fmt.Sprintf(":%d", config.Http.Port), Middleware(mux))
}

func Middleware(h http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println(request.URL)
		h.ServeHTTP(writer, request)
	}
}
