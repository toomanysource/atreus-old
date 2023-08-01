// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"Atreus/app/relation/service/internal/biz"
	"Atreus/app/relation/service/internal/conf"
	"Atreus/app/relation/service/internal/data"
	"Atreus/app/relation/service/internal/server"
	"Atreus/app/relation/service/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, logger log.Logger) (*kratos.App, func(), error) {
	dataData, cleanup, err := data.NewData(logger)
	if err != nil {
		return nil, nil, err
	}
	relationRepo := data.NewRelationRepo(dataData, logger)
	relationUsecase := biz.NewRelationUsecase(relationRepo, logger)
	relationService := service.NewRelationService(relationUsecase, logger)
	grpcServer := server.NewGRPCServer(confServer, relationService, logger)
	httpServer := server.NewHTTPServer(confServer, relationService, logger)
	app := newApp(logger, grpcServer, httpServer)
	return app, func() {
		cleanup()
	}, nil
}
