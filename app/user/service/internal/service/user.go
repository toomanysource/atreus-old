package service

import (
	"Atreus/app/user/service/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/log"

	pb "Atreus/api/user/service/v1"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	log *log.Helper

	usecase *biz.UserUsecase
}

func NewUserService(uc *biz.UserUsecase, logger log.Logger) *UserService {
	return &UserService{usecase: uc, log: log.NewHelper(logger)}
}

func (s *UserService) UserRegister(ctx context.Context, req *pb.UserRegisterRequest) (*pb.UserRegisterReply, error) {
	user, err := s.usecase.Register(ctx, req.Username, req.Password)
	if err != nil {
		return &pb.UserRegisterReply{
			StatusCode: 300,
			StatusMsg:  "failed",
		}, nil
	}

	return &pb.UserRegisterReply{
		StatusCode: 200,
		StatusMsg:  "success",
		UserId:     user.ID,
	}, nil
}

func (s *UserService) UserLogin(ctx context.Context, req *pb.UserLoginRequest) (*pb.UserLoginReply, error) {
	user, err := s.usecase.Login(ctx, req.Username, req.Password)
	if err != nil {
		return &pb.UserLoginReply{
			StatusCode: 300,
			StatusMsg:  "failed",
		}, nil
	}

	return &pb.UserLoginReply{
		StatusCode: 200,
		StatusMsg:  "success",
		UserId:     user.ID,
	}, nil
}

func (s *UserService) GetUserInfo(ctx context.Context, req *pb.UserInfoRequest) (*pb.UserInfoReply, error) {
	user, err := s.usecase.GetInfo(ctx, req.UserId)
	if err != nil {
		return &pb.UserInfoReply{
			StatusCode: 300,
			StatusMsg:  "failed",
		}, nil
	}

	return &pb.UserInfoReply{
		StatusCode: 200,
		StatusMsg:  "success",
		User: &pb.User{
			Id:              user.ID,
			Name:            user.Name,
			FollowCount:     user.FollowCount,
			FollowerCount:   user.FollowerCount,
			IsFollow:        user.IsFollow,
			Avatar:          user.Avatar,
			BackgroundImage: user.BackgroundImage,
			Signature:       user.Signature,
			TotalFavorite:   user.TotalFavorite,
			WorkCount:       user.WorkCount,
			FavoriteCount:   user.FavoriteCount,
		},
	}, nil
}
