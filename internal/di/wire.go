// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package di

import (
	"github.com/cacos-group/cacos/internal/conf"
	"github.com/cacos-group/cacos/internal/core/event/sourcing"
	"github.com/cacos-group/cacos/internal/core/query"
	"github.com/cacos-group/cacos/internal/logger"
	"github.com/cacos-group/cacos/internal/resource"
	"github.com/cacos-group/cacos/internal/server/grpc"
	"github.com/cacos-group/cacos/internal/service"
	"github.com/google/wire"
)

//bash ~/go/bin/wire
func InitApp(config *conf.Config) (*App, func(), error) {

	panic(wire.Build(NewApp, service.Provider, grpc.New, logger.NewLog, query.NewClient, sourcing.NewClient, resource.NewDB, resource.NewEtcd))

}
