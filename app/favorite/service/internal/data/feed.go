package data

import (
	pb "Atreus/api/feed/service/v1"
	"google.golang.org/grpc"
)

type FeedRepo struct {
	client pb.FeedServiceClient
}

func NewFeedRepo(conn *grpc.ClientConn) *FeedRepo {
	return &FeedRepo{
		client: pb.NewFeedServiceClient(conn),
	}
}

// GetVideoInfoByVideoIds 通过videoId获取视频信息; TODO: to be implemented by feed service
//func (f *FeedRepo) GetVideoInfoByVideoIds(
//	ctx context.Context, videoIds []uint32) ([]biz.Video, error) {
//	// call grpc function to fetch video info
//	resp, err := f.client.GetVideoInfoByVideoIds(ctx, &pb.ClientVideoInfoRequest{VideoIds: videoIds})
//	if err != nil {
//		return nil, err
//	}
//	if len(resp.videos) == 0 {
//		return nil, errors.New("video not found")
//	}
//	// convert pb.Video slice to biz.Video slice
//	videos := make([]biz.Video, len(resp.videos))
//	if err := copier.Copy(&videos, &resp.videos); err != nil {
//		return nil, err
//	}
//	return videos, nil
//}
