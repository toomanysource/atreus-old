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

//func (c *UserServiceClient) GetUserInfoByUserIds(ctx context.Context, userId uint32) (*pb.ClientUserInfoReply, error) {
//	return c.uc.GetUserInfoByUserIds(ctx, &pb.ClientUserInfoByUserIdRequest{UserId: userId})
//}

func (c *UserServiceClient) GetUserInfoByUserIds(
	ctx context.Context, userIds []uint32) (*pb.ClientUserInfosReply, error) {
	return c.uc.GetUserInfoByUserIds(ctx, &pb.ClientUserInfoByUserIdsRequest{UserId: userIds})
}
