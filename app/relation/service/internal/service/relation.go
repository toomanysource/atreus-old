package service

import (
	"context"

	pb "Atreus/api/relation/service/v1"
)

type RelationService struct {
	pb.UnimplementedRelationServiceServer
}

func NewRelationService() *RelationService {
	return &RelationService{}
}

func (s *RelationService) GetFollowerRelationList(ctx context.Context, req *pb.RelationFollowerListRequest) (*pb.RelationFollowerListReply, error) {
	return &pb.RelationFollowerListReply{}, nil
}
func (s *RelationService) GetFollowRelationList(ctx context.Context, req *pb.RelationFollowListRequest) (*pb.RelationFollowListReply, error) {
	return &pb.RelationFollowListReply{}, nil
}
func (s *RelationService) RelationAction(ctx context.Context, req *pb.RelationActionRequest) (*pb.RelationActionReply, error) {
	return &pb.RelationActionReply{}, nil
}
