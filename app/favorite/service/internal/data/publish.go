package data

import (
	pb "Atreus/api/publish/service/v1"
	"Atreus/app/favorite/service/internal/biz"
	"Atreus/app/favorite/service/internal/server"
	"context"
	"errors"
	"github.com/jinzhu/copier"
)

type publishRepo struct {
	client pb.PublishServiceClient
}

func NewPublishRepo(conn server.UserConn) biz.PublishRepo {
	return &publishRepo{
		client: pb.NewPublishServiceClient(conn),
	}
}

// GetVideoListByVideoIds 通过videoId获取视频信息;
func (f *publishRepo) GetVideoListByVideoIds(
	ctx context.Context, videoIds []uint32) ([]biz.Video, error) {
	// call grpc function to fetch video info
	resp, err := f.client.GetVideoListByVideoIds(ctx, &pb.VideoListByVideoIdsRequest{VideoId: videoIds})
	if err != nil {
		return nil, err
	}
	if len(resp.VideoList) == 0 {
		return nil, errors.New("video not found")
	}
	// convert pb.Video slice to biz.Video slice
	videos := make([]biz.Video, len(resp.VideoList))
	if err := copier.Copy(&videos, &resp.VideoList); err != nil {
		return nil, err
	}
	return videos, nil
}
