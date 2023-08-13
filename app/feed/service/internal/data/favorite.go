package data

import (
	pb "Atreus/api/favorite/service/v1"
	"Atreus/app/feed/service/internal/server"
	"context"
)

type favoriteRepo struct {
	client pb.FavoriteServiceClient
}

func NewFavoriteRepo(conn server.FavoriteConn) FavoriteRepo {
	return &favoriteRepo{
		client: pb.NewFavoriteServiceClient(conn),
	}
}

func (u *favoriteRepo) IsFavorite(ctx context.Context, userId uint32, videoId uint32) (bool, error) {
	resp, err := u.client.IsFavorite(ctx, &pb.IsFavoriteRequest{UserId: userId, VideoId: videoId})
	if err != nil {
		return false, err
	}
	result := resp.GetIsFavorite()
	return result, nil
}
