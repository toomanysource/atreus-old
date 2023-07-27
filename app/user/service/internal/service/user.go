package service

import (
	"context"

	pb "Atreus/api/user/service/v1"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) UserRegister(ctx context.Context, req *pb.UserRegisterRequest) (*pb.UserRegisterReply, error) {
	return &pb.UserRegisterReply{}, nil
}
func (s *UserService) UserLogin(ctx context.Context, req *pb.UserLoginRequest) (*pb.UserLoginReply, error) {
	return &pb.UserLoginReply{}, nil
}
func (s *UserService) GetUserInfo(ctx context.Context, req *pb.UserInfoRequest) (*pb.UserInfoReply, error) {
	return &pb.UserInfoReply{}, nil
}
