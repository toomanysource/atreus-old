package service

import (
	"context"

	pb "Atreus/api/favorite/service/v1"
)

type FavoriteService struct {
	pb.UnimplementedFavoriteServiceServer
}

func NewFavoriteService() *FavoriteService {
	return &FavoriteService{}
}

func (s *FavoriteService) GetFavoriteList(ctx context.Context, req *pb.FavoriteListRequest) (*pb.FavoriteListReply, error) {
	return &pb.FavoriteListReply{}, nil
}
func (s *FavoriteService) FavoriteAction(ctx context.Context, req *pb.FavoriteActionRequest) (*pb.FavoriteActionReply, error) {
	return &pb.FavoriteActionReply{}, nil
}
