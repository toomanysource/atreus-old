package service

import (
	"context"

	pb "Atreus/api/publish/service/v1"
)

type PublishService struct {
	pb.UnimplementedPublishServer
}

func NewPublishService() *PublishService {
	return &PublishService{}
}

func (s *PublishService) GetPublishList(ctx context.Context, req *pb.PublishListRequest) (*pb.PublishListReply, error) {
	return &pb.PublishListReply{}, nil
}
func (s *PublishService) PublishAction(ctx context.Context, req *pb.PublishActionRequest) (*pb.PublishActionReplay, error) {
	return &pb.PublishActionReplay{}, nil
}
