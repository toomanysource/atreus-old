package service

import (
	"context"

	// v1 "Atreus/api/feed/service/v1"
	pb "Atreus/api/feed/service/v1"

	"Atreus/app/feed/service/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type FeedService struct {
	pb.UnimplementedFeedServiceServer
	log *log.Helper
	fu  *biz.FeedUsecase
}

func NewFeedService(fu *biz.FeedUsecase, logger log.Logger) *FeedService {
	return &FeedService{
		fu:  fu,
		log: log.NewHelper(log.With(logger, "model", "service/feed")),
	}
}

// FeedList 返回一个按照投稿时间倒序的视频列表，单次最多30个视频
func (s *FeedService) FeedList(ctx context.Context, req *pb.ListFeedRequest) (*pb.ListFeedReply, error) {
	var nextTime int64
	reply := &pb.ListFeedReply{StatusCode: 0, StatusMsg: "Success", VideoList: make([]*pb.Video, 0), NextTime: 0}
	// videos, nextTime, err := s.fu.GetFeedList(ctx, req.LatestTime)
	nextTime, videos, err := s.fu.Feedlist(ctx, req.LatestTime, req.Token)
	if err != nil {
		reply.StatusCode = -1
		reply.StatusMsg = err.Error()
		return reply, nil
	}

	for _, video := range videos {
		reply.VideoList = append(reply.VideoList, &pb.Video{
			Id:    video.Id,
			Title: video.Title,
			Author: &pb.User{
				Id:              video.Author.Id,
				Name:            video.Author.Name,
				Avatar:          video.Author.Avatar,
				BackgroundImage: video.Author.BackgroundImage,
				Signature:       video.Author.Signature,
				IsFollow:        video.Author.IsFollow,
				FollowCount:     video.Author.FollowCount,
				FollowerCount:   video.Author.FollowerCount,
				TotalFavorited:  video.Author.TotalFavorited,
				WorkCount:       video.Author.WorkCount,
				FavoriteCount:   video.Author.FavoriteCount,
			},
			PlayUrl:       video.PlayURL,
			CoverUrl:      video.CoverURL,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    video.IsFavorite,
		})
	}
	reply.NextTime = nextTime

	return reply, nil
}
