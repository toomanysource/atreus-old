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

var ProviderSet = wire.NewSet(NewGRPCServer, NewHTTPServer, NewFeedClient)

func NewFeedClient(c *conf.Client, logger log.Logger) *stdgrpc.ClientConn {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(c.User.To),
		grpc.WithMiddleware(
			recovery.Recovery(),
			logging.Server(logger),
		),
	)
	if err != nil {
		log.Fatalf("Error when Favorite connecting to Feed Services, err : %w", err)
	}
	return conn
}
