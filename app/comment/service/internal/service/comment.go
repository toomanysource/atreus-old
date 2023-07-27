package service

import (
	"context"

	pb "Atreus/api/comment/service/v1"
)

type CommentService struct {
	pb.UnimplementedCommentServiceServer
}

func NewCommentService() *CommentService {
	return &CommentService{}
}

func (s *CommentService) GetCommentList(ctx context.Context, req *pb.CommentListRequest) (*pb.CommentListReply, error) {
	return &pb.CommentListReply{}, nil
}
func (s *CommentService) CommentAction(ctx context.Context, req *pb.CommentActionRequest) (*pb.CommentActionReply, error) {
	return &pb.CommentActionReply{}, nil
}
