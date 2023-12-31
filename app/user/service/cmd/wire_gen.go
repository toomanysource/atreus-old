// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/toomanysource/atreus/app/user/service/internal/biz"
	"github.com/toomanysource/atreus/app/user/service/internal/conf"
	"github.com/toomanysource/atreus/app/user/service/internal/data"
	"github.com/toomanysource/atreus/app/user/service/internal/server"
	"github.com/toomanysource/atreus/app/user/service/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"

	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, confData *conf.Data, client *conf.Client, jwt *conf.JWT, logger log.Logger) (*kratos.App, func(), error) {
	db := data.NewGormDb(confData)
	dataData, cleanup, err := data.NewData(db, logger)
	if err != nil {
		return nil, nil, err
	}
	relationConn := server.NewRelationClient(client, logger)
	userRepo := data.NewUserRepo(dataData, relationConn, logger)
	userUsecase := biz.NewUserUsecase(userRepo, jwt, logger)
	userService := service.NewUserService(userUsecase, logger)
	grpcServer := server.NewGRPCServer(confServer, userService, logger)
	httpServer := server.NewHTTPServer(confServer, jwt, userService, logger)
	app := newApp(logger, grpcServer, httpServer)
	return app, func() {
		cleanup()
	}, nil
}
