package biz

import (
	"Atreus/app/favorite/service/internal/conf"
	"Atreus/pkg/common"
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

// Favorite is corresponding to the favorite table in database
type Favorite struct {
	//ID      int -> use <Composite Primary Key> instead
	VideoID   uint32
	UserID    uint32
	CreatedAt time.Time
	DeletedAt *time.Time
}

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

// FavoriteRepo is database manipulation interface
type FavoriteRepo interface {
	GetFavoriteList(ctx context.Context, userID uint32) ([]Video, error)             // list of user's favorite video; use slice without pointer
	IsFavorite(ctx context.Context, userID uint32, videoID []uint32) ([]bool, error) // whether a list of video is favorited by a user
	DeleteFavoriteTx(ctx context.Context, userID uint32, videoID uint32) error
	CreateFavoriteTx(ctx context.Context, userID uint32, videoID uint32) error
}

type UserRepo interface {
	UpdateFavorited(ctx context.Context, userId uint32, change int32) error
	UpdateFavorite(ctx context.Context, userId uint32, change int32) error
}

type PublishRepo interface {
	GetVideoListByVideoIds(ctx context.Context, videoIds []uint32) ([]Video, error) // 多个/单个视频信息
	UpdateFavoriteCount(ctx context.Context, videoId uint32, change int32) error
}

// FavoriteUsecase .
type FavoriteUsecase struct {
	favoriteRepo FavoriteRepo
	config       *conf.JWT
	log          *log.Helper
}

// NewFavoriteUsecase clear unnecessary dependencies
func NewFavoriteUsecase(conf *conf.JWT, repo FavoriteRepo, logger log.Logger) *FavoriteUsecase {
	return &FavoriteUsecase{config: conf, favoriteRepo: repo, log: log.NewHelper(log.With(logger, "model", "usecase/favorite"))}
}

// FavoriteAction is for http api use; create & delete integrated
func (uc *FavoriteUsecase) FavoriteAction(ctx context.Context, videoId, actionType uint32, tokenString string) error {
	// user verification & get user_id
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
	// biz
	switch actionType {
	case 1:
		return uc.favoriteRepo.CreateFavoriteTx(ctx, userId, videoId)
	case 2:
		return uc.favoriteRepo.DeleteFavoriteTx(ctx, userId, videoId)
	default:
		return errors.New("invalid action type(not 1 nor 2)")
	}
}

func (uc *FavoriteUsecase) GetFavoriteList(ctx context.Context, userID uint32, tokenString string) ([]Video, error) {
	// user verification & get user_id
	token, err := common.ParseToken(uc.config.Http.TokenKey, tokenString)
	if err != nil {
		return nil, err
	}
	tokenData, err := common.GetTokenData(token)
	if err != nil {
		return nil, err
	}
	userIDFloat64, ok := tokenData["user_id"].(float64)
	if !ok {
		return nil, errors.New("user_id is not a valid float64")
	}

	userIDFromToken := uint32(userIDFloat64)
	if userIDFromToken != userID {
		return nil, errors.New("user_id and token not match")
	}
	// biz
	return uc.favoriteRepo.GetFavoriteList(ctx, userID)
	//return nil, nil
}

func (uc *FavoriteUsecase) IsFavorite(ctx context.Context, userID uint32, videoIDs []uint32) ([]bool, error) {
	// internal use; no need to verify token
	ret, err := uc.favoriteRepo.IsFavorite(ctx, userID, videoIDs)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
