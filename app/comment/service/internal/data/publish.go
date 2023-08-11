package data

import (
	pb "Atreus/api/publish/service/v1"
	"Atreus/app/comment/service/internal/server"
	"context"
	"errors"
)

type publishRepo struct {
	client pb.PublishServiceClient
}

func NewPublishRepo(conn server.PublishConn) PublishRepo {
	return &publishRepo{
		client: pb.NewPublishServiceClient(conn),
	}
}

// UpdateComment 接收Publish服务的回应
func (u *publishRepo) UpdateComment(ctx context.Context, videoId uint32, commentChange int32) error {
	resp, err := u.client.UpdateComment(
		ctx, &pb.UpdateCommentCountRequest{VideoId: videoId, CommentChange: commentChange})
	if err != nil {
		return err
	}
	if resp.StatusCode != 0 {
		return errors.New(resp.StatusMsg)
	}
	return nil
}
