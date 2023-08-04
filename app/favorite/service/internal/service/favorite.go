package service

import (
	pb "Atreus/api/favorite/service/v1"
	"Atreus/app/favorite/service/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
)

type FavoriteService struct {
	pb.UnimplementedFavoriteServiceServer
	fu  *biz.FavoriteUsecase
	log *log.Helper
}

func NewFavoriteService(fu *biz.FavoriteUsecase, logger log.Logger) *FavoriteService {
	return &FavoriteService{
		fu:  fu,
		log: log.NewHelper(log.With(logger, "model", "service/favorite")),
	}
}
