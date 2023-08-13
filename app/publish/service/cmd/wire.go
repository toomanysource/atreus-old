//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"Atreus/app/publish/service/internal/biz"
	"Atreus/app/publish/service/internal/conf"
	"Atreus/app/publish/service/internal/data"
	"Atreus/app/publish/service/internal/server"
	"Atreus/app/publish/service/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Client, *conf.Minio, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
