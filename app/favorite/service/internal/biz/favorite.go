package biz

import (
	"Atreus/app/favorite/service/internal/conf"
	"Atreus/pkg/common"
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
)

// Video is used to receive video info from video service;response is an array of Videos
type Video struct {
	Id            uint32
	Author        *User
	PlayUrl       string
	CoverUrl      string
	FavoriteCount uint32
	CommentCount  uint32
	IsFavorite    bool
	Title         string
}

// User is used to receive video info from user service;& send response
type User struct {
	Id              uint32
	Name            string
	FollowCount     uint32
	FollowerCount   uint32
	IsFollow        bool
	Avatar          string
	BackgroundImage string
	Signature       string
	TotalFavorited  uint32 // 总获赞数
	WorkCount       uint32
	FavoriteCount   uint32 // 点赞数量
}

type FavoriteRepo interface {
	GetFavoriteList(ctx context.Context, userID uint32) ([]Video, error)
	IsFavorite(ctx context.Context, userID uint32, videoID []uint32) ([]bool, error)
	DelFavorite(ctx context.Context, userID uint32, videoID uint32) error
	AddFavorite(ctx context.Context, userID uint32, videoID uint32) error
}

type UserRepo interface {
	UpdateFavorited(ctx context.Context, userId uint32, change int32) error
	UpdateFavorite(ctx context.Context, userId uint32, change int32) error
}

type PublishRepo interface {
	GetVideoListByVideoIds(ctx context.Context, userId uint32, videoIds []uint32) ([]Video, error)
	UpdateFavoriteCount(ctx context.Context, videoId uint32, change int32) error
}

type FavoriteUsecase struct {
	favoriteRepo FavoriteRepo
	config       *conf.JWT
	log          *log.Helper
}

func NewFavoriteUsecase(conf *conf.JWT, repo FavoriteRepo, logger log.Logger) *FavoriteUsecase {
	return &FavoriteUsecase{config: conf, favoriteRepo: repo, log: log.NewHelper(log.With(logger, "model", "usecase/favorite"))}
}

func (uc *FavoriteUsecase) FavoriteAction(ctx context.Context, videoId, actionType uint32, tokenString string) error {

	token, err := common.ParseToken(uc.config.Http.TokenKey, tokenString)
	if err != nil {
		return err
	}
	tokenData, err := common.GetTokenData(token)
	if err != nil {
		return err
	}
	userIDFloat64, ok := tokenData["user_id"].(float64)
	if !ok {
		return errors.New("user_id is not a valid float64")
	}
	userId := uint32(userIDFloat64)

	switch actionType {
	case 1:
		return uc.favoriteRepo.AddFavorite(ctx, userId, videoId)
	case 2:
		return uc.favoriteRepo.DelFavorite(ctx, userId, videoId)
	default:
		return errors.New("invalid action type(not 1 nor 2)")
	}
}

func (uc *FavoriteUsecase) GetFavoriteList(ctx context.Context, userID uint32, tokenString string) ([]Video, error) {
	if tokenString != "" {
		token, err := common.ParseToken(uc.config.Http.TokenKey, tokenString)
		if err != nil {
			return nil, err
		}
		_, err = common.GetTokenData(token)
		if err != nil {
			return nil, err
		}
	}
	return uc.favoriteRepo.GetFavoriteList(ctx, userID)
}

func (uc *FavoriteUsecase) IsFavorite(ctx context.Context, userID uint32, videoIDs []uint32) ([]bool, error) {

	ret, err := uc.favoriteRepo.IsFavorite(ctx, userID, videoIDs)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
