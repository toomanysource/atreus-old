package data

import (
	pb "Atreus/api/user/service/v1"
	"Atreus/app/relation/service/internal/biz"
	"context"
	"errors"
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

// GetUserInfoByUserIds 接收User服务的回应，并转化为biz.User类型
func (u *UserRepo) GetUserInfoByUserIds(ctx context.Context, userIds []uint32) ([]*biz.User, error) {
	resp, err := u.client.GetUserInfoByUserIds(ctx, &pb.ClientUserInfoByUserIdsRequest{UserIds: userIds})
	if err != nil {
		return nil, err
	}

	// 判空
	if len(resp.Users) == 0 {
		return nil, errors.New("the user service did not search for any information")
	}

	users := make([]*biz.User, 0, len(resp.Users)+1)
	for _, user := range resp.Users {
		users = append(users, &biz.User{
			Id:              user.Id,
			Name:            user.Name,
			Avatar:          user.Avatar,
			BackgroundImage: user.BackgroundImage,
			Signature:       user.Signature,
			IsFollow:        user.IsFollow,
			FollowCount:     user.FollowCount,
			FollowerCount:   user.FollowerCount,
			TotalFavorite:   user.TotalFavorited,
			WorkCount:       user.WorkCount,
			FavoriteCount:   user.FavoriteCount,
		})
	}
	return users, nil
}
