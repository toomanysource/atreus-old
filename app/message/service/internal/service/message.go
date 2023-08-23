package service

import (
	pb "Atreus/api/message/service/v1"
	"Atreus/app/message/service/internal/biz"
	"context"
	"github.com/jinzhu/copier"

	"github.com/go-kratos/kratos/v2/log"
)

type MessageService struct {
	pb.UnimplementedMessageServiceServer
	mu  *biz.MessageUsecase
	log *log.Helper
}

func NewMessageService(mu *biz.MessageUsecase, logger log.Logger) *MessageService {
	return &MessageService{
		mu:  mu,
		log: log.NewHelper(log.With(logger, "model", "service/Message")),
	}
}

func (s *MessageService) GetMessageList(ctx context.Context, req *pb.MessageListRequest) (*pb.MessageListReply, error) {
	message, err := s.mu.GetMessageList(ctx, req.Token, req.ToUserId, req.PreMsgTime)
	if err != nil {
		return &pb.MessageListReply{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		}, nil
	}
	ml := make([]*pb.Message, len(message))
	if err = copier.Copy(&ml, message); err != nil {
		return &pb.MessageListReply{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		}, nil
	}
	return &pb.MessageListReply{
		StatusCode:  0,
		StatusMsg:   "success",
		MessageList: ml,
	}, nil
}
func (s *MessageService) MessageAction(ctx context.Context, req *pb.MessageActionRequest) (*pb.MessageActionReply, error) {
	reply := &pb.MessageActionReply{StatusCode: 0, StatusMsg: "success"}
	err := s.mu.PublishMessage(ctx, req.Token, req.ToUserId, req.ActionType, req.Content)
	if err != nil {
		reply.StatusCode = -1
		reply.StatusMsg = err.Error()
		return reply, nil
	}
	return reply, nil
}
