package service

import (
	pb "Atreus/api/feed/service/v1"
	"google.golang.org/grpc"
)

// FeedServiceClient 设置一个User服务客户端实体
type FeedServiceClient struct {
	client pb.FeedServiceClient
	conn   *grpc.ClientConn
}

func NewFeedServiceClient(conn *grpc.ClientConn) *FeedServiceClient {
	return &FeedServiceClient{
		conn:   conn,
		client: pb.NewFeedServiceClient(conn),
	}
}

// GetVideoInfo 通过videoId获取视频信息；TODO: to be implemented by feed service
//func (c *FeedServiceClient) GetVideoInfoByVideoIds(
//	ctx context.Context, videoIds []uint32) (*pb.ClientVideoInfoReply, error) {
//	return c.client.GetVideoInfoByVideoIds(ctx, &pb.ClientVideoInfoRequest{VideoIds: videoIds})
//}
