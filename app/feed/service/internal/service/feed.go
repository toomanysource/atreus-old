package service

import (
	"context"

	pb "Atreus/api/feed/service/v1"
)

type FeedService struct {
	pb.UnimplementedFeedServer
}

func NewFeedService() *FeedService {
	return &FeedService{}
}

func (s *FeedService) GetFeed(ctx context.Context, req *pb.FeedRequest) (*pb.FeedReply, error) {
	return &pb.FeedReply{}, nil
}
