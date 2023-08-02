package service

import (
	pb "Atreus/api/relation/service/v1"
	"Atreus/app/relation/service/internal/biz"
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type RelationService struct {
	pb.UnimplementedRelationServiceServer
	log     *log.Helper
	usecase *biz.RelationUsecase
}

func NewRelationService(uc *biz.RelationUsecase, logger log.Logger) *RelationService {
	return &RelationService{usecase: uc, log: log.NewHelper(logger)}
}

// RelationAction 关注/取消关注
func (s *RelationService) RelationAction(ctx context.Context, req *pb.RelationActionRequest) (*pb.RelationActionReply, error) {
	err := s.usecase.Action(ctx, req.Token, req.ToUserId, req.ActionType)
	if err != nil {
		return &pb.RelationActionReply{
			StatusCode: 300,
			StatusMsg:  "failed",
		}, nil
	}
	return &pb.RelationActionReply{
		StatusCode: 200,
		StatusMsg:  "success",
	}, nil
}

// GetFollowRelationList 获取关注列表
func (s *RelationService) GetFollowRelationList(ctx context.Context, req *pb.RelationFollowListRequest) (*pb.RelationFollowListReply, error) {
	reply := &pb.RelationFollowListReply{StatusCode: 200, StatusMsg: "Success", UserList: make([]*pb.User, 0)}
	//list, err := s.usecase.GetFollowList(ctx, req.UserId, req.Token)
	//if err != nil {
	//	return &pb.RelationFollowListReply{
	//		StatusCode: 300,
	//		StatusMsg:  "failed",
	//	}, nil
	//}
	return reply, nil
}

// GetFollowerRelationList 获取粉丝列表
func (s *RelationService) GetFollowerRelationList(ctx context.Context, req *pb.RelationFollowerListRequest) (*pb.RelationFollowerListReply, error) {
	reply := &pb.RelationFollowerListReply{StatusCode: 200, StatusMsg: "Success", UserList: make([]*pb.User, 0)}
	//list, err := s.usecase.GetFollowerList(ctx, req.UserId, req.Token)
	//if err != nil {
	//	return &pb.RelationFollowerListReply{
	//		StatusCode: 300,
	//		StatusMsg:  "failed",
	//	}, nil
	//}
	return reply, nil
}
