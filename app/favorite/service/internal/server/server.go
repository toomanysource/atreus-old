package server

import (
	"Atreus/app/favorite/service/internal/conf"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
	stdgrpc "google.golang.org/grpc"
)

var ProviderSet = wire.NewSet(NewGRPCServer, NewHTTPServer, NewPublishClient, NewUserClient)

type UserConn stdgrpc.ClientConnInterface
type PublishConn stdgrpc.ClientConnInterface

// NewUserClient 创建一个User服务客户端，接收User服务数据
func NewUserClient(c *conf.Client, logger log.Logger) UserConn {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(c.User.To),
		grpc.WithMiddleware(
			recovery.Recovery(),
			logging.Client(logger),
		),
	)
	if err != nil {
		log.Fatalf("Error connecting to User Services, err : %w", err)
	}
	return conn
}

// NewPublishClient 创建一个Publish服务客户端，接收Publish服务数据
func NewPublishClient(c *conf.Client, logger log.Logger) PublishConn {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(c.Publish.To),
		grpc.WithMiddleware(
			recovery.Recovery(),
			logging.Client(logger),
		),
	)
	if err != nil {
		log.Fatalf("Error connecting to Publish Services, err : %w", err)
	}
	return conn
}
