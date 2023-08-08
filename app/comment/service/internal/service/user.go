package service

import (
	pb "Atreus/api/user/service/v1"
	"context"
	"google.golang.org/grpc"
)

// UserServiceClient 设置一个User服务客户端实体
type UserServiceClient struct {
	uc   pb.UserServiceClient
	conn *grpc.ClientConn
}

func NewUserServiceClient(conn *grpc.ClientConn) *UserServiceClient {
	return &UserServiceClient{
		conn: conn,
		uc:   pb.NewUserServiceClient(conn),
	}
}

func (c *UserServiceClient) GetUserInfos(
	ctx context.Context, userIds []uint32) (*pb.UserInfosReply, error) {
	return c.uc.GetUserInfos(ctx, &pb.UserInfosRequest{UserIds: userIds})
}