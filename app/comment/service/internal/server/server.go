package server

import (
	"Atreus/app/comment/service/internal/conf"
	"context"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	jwtv4 "github.com/golang-jwt/jwt/v4"
	"github.com/google/wire"
	stdgrpc "google.golang.org/grpc"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewGRPCServer, NewHTTPServer, NewUserClient)

// NewUserClient 创建一个User服务客户端，接收User服务数据
func NewUserClient(c *conf.Client, j *conf.JWT) *stdgrpc.ClientConn {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(c.User.To),
		grpc.WithMiddleware(jwt.Client(func(token *jwtv4.Token) (interface{}, error) {
			return []byte(j.Grpc.TokenKey), nil
		})),
	)
	if err != nil {
		panic(err)
	}
	return conn
}
