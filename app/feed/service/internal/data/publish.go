package data

import (
	pb "Atreus/api/publish/service/v1"
	"Atreus/app/feed/service/internal/biz"
	"Atreus/app/feed/service/internal/server"
	"context"
	"errors"

	"github.com/jinzhu/copier"
)

type publishRepo struct {
	client pb.PublishServiceClient
}

func NewPublishRepo(conn server.PublishConn) biz.PublishRepo {
	return &publishRepo{
		client: pb.NewPublishServiceClient(conn),
	}
}

var (
	videoCount uint32 = 30
)

// GetVideoList
func (u *publishRepo) GetVideoList(
	ctx context.Context, latestTime string, userId uint32, number uint32) (int64, []biz.Video, error) {

	// call grpc function to fetch video list
	resp, _ := u.client.GetVideoList(ctx, &pb.VideoListRequest{LatestTime: latestTime, UserId: userId, Number: videoCount})
	if resp.StatusCode != 0 {
		return 0, nil, errors.New(resp.StatusMsg)
	}

	// convert pb.Video slice to biz.Video slice
	videos := make([]biz.Video, len(resp.VideoList))
	if err := copier.Copy(&videos, &resp.VideoList); err != nil {
		return 0, nil, err
	}

	nextTime := resp.NextTime
	return nextTime, videos, nil
}
