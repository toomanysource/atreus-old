package data

import (
	pb "Atreus/api/user/service/v1"
	"Atreus/app/comment/service/internal/biz"
	"context"
	"google.golang.org/grpc"
)

type UserRepo struct {
	client pb.UserServiceClient
}

func NewUserRepo(conn *grpc.ClientConn) *UserRepo {
	return &UserRepo{
		client: pb.NewUserServiceClient(conn),
	}
}

// GetUserInfoByUserId 接收User服务的回应，并转化为biz.User类型
func (u *UserRepo) GetUserInfoByUserId(ctx context.Context, userId uint32) (biz.User, error) {
	resp, err := u.client.GetUserInfoByUserId(ctx, &pb.ClientUserInfoByUserIdRequest{UserId: userId})
	if err != nil {
		return biz.User{}, err
	}
	return biz.User{
		Id:              resp.User.Id,
		Name:            resp.User.Name,
		Avatar:          resp.User.Avatar,
		BackgroundImage: resp.User.BackgroundImage,
		Signature:       resp.User.Signature,
		IsFollow:        resp.User.IsFollow,
		FollowCount:     resp.User.FollowCount,
		FollowerCount:   resp.User.FollowerCount,
		TotalFavorited:  resp.User.TotalFavorited,
		WorkCount:       resp.User.WorkCount,
		FavoriteCount:   resp.User.FavoriteCount,
	}, nil
}
