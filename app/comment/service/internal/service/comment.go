package service

import (
	"context"

	pb "Atreus/api/comment/service/v1"
	"Atreus/app/comment/service/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type CommentService struct {
	pb.UnimplementedCommentServiceServer
	cu  *biz.CommentUsecase
	log *log.Helper
}

func NewCommentService(cu *biz.CommentUsecase, logger log.Logger) *CommentService {
	return &CommentService{
		cu:  cu,
		log: log.NewHelper(log.With(logger, "model", "service/comment")),
	}
}

func (s *CommentService) GetCommentList(ctx context.Context, req *pb.CommentListRequest) (*pb.CommentListReply, error) {
	return &pb.CommentListReply{
		StatusCode:  0,
		StatusMsg:   "Success",
		CommentList: []*pb.Comment{},
	}, nil
}
func (s *CommentService) CommentAction(ctx context.Context, req *pb.CommentActionRequest) (*pb.CommentActionReply, error) {
	return &pb.CommentActionReply{
		StatusCode: 0,
		StatusMsg:  "Success",
		Comment:    &pb.Comment{},
	}, nil
}
