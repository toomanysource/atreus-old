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
	TotalFavorited  uint32
	WorkCount       uint32
	FavoriteCount   uint32
}

// FavoriteRepo is database manipulation interface
type FavoriteRepo interface {
	CreateFavorite(context context.Context, userId, videoId uint32) error
	DeleteFavorite(context context.Context, userId, videoId uint32) error
	GetFavoriteList(context context.Context, userID uint32) ([]Video, error)  // list of user's favorite video; use slice without pointer
	IsFavorite(context context.Context, userID, videoID uint32) (bool, error) // whether a video is favorited by a user
	// deprecated
	//CountFavoriteByVideoIDs(context context.Context, videoIDs []uint32) (int64, error) // num of favorite of a video
	//CountFavoriteByUserID(context context.Context, userID uint32) (int64, error)       // num of favorite user received
}

// FavoriteUsecase .
type FavoriteUsecase struct {
	favoriteRepo FavoriteRepo
	userRepo     UserRepo
	publishRepo  PublishRepo
	tx           Transaction // transaction is used to support consistency
	config       *conf.JWT
	log          *log.Helper
}

type UserRepo interface {
	UpdateFavorited(ctx context.Context, userId uint32, change int32) error
	UpdateFavorite(ctx context.Context, userId uint32, change int32) error
}

type PublishRepo interface {
	GetVideoListByVideoIds(ctx context.Context, videoIds []uint32) ([]Video, error) // 多个/单个视频信息
}

func NewFavoriteUsecase(conf *conf.JWT, repo FavoriteRepo, ur UserRepo, pr PublishRepo, trans Transaction, logger log.Logger) *FavoriteUsecase {
	return &FavoriteUsecase{config: conf, favoriteRepo: repo, userRepo: ur, publishRepo: pr, tx: trans, log: log.NewHelper(logger)}
}

// FavoriteAction is for http api use; create & delete integrated
func (uc *FavoriteUsecase) FavoriteAction(ctx context.Context, videoId, actionType uint32, tokenString string) error {
	// user verification & get user_id
	token, err := common.ParseToken(uc.config.Http.TokenKey, tokenString)
	if err != nil {
		return err
	}
	data, err := common.GetTokenData(token)
	if err != nil {
		return err
	}
	userIDFloat64, ok := data["user_id"].(float64)
	if !ok {
		return errors.New("user_id is not a valid float64")
	}
	userId := uint32(userIDFloat64)
	// biz
	videoInfo, err := uc.publishRepo.GetVideoListByVideoIds(ctx, []uint32{videoId})
	if err != nil {
		return err
	}
	if len(videoInfo) == 0 {
		return errors.New("video not found")
	}
	authorId := videoInfo[0].Author.Id // to update author's favorited
	switch actionType {
	case 1:
		return uc.tx.ExecTx(ctx, func(ctx context.Context) error {
			// create favorite
			err2 := uc.favoriteRepo.CreateFavorite(ctx, userId, videoId)
			if err2 != nil {
				return err2
			}
			// notify other services
			err2 = uc.userRepo.UpdateFavorited(ctx, authorId, 1)
			if err2 != nil {
				return err2
			}
			err2 = uc.userRepo.UpdateFavorite(ctx, userId, 1)
			if err2 != nil {
				return err2
			}
			return nil
		})
	case 2:
		return uc.tx.ExecTx(ctx, func(ctx context.Context) error {
			// delete favorite
			err2 := uc.favoriteRepo.DeleteFavorite(ctx, userId, videoId)
			if err2 != nil {
				return err2
			}
			// notify other services
			err2 = uc.userRepo.UpdateFavorited(ctx, authorId, -1)
			if err2 != nil {
				return err2
			}
			err2 = uc.userRepo.UpdateFavorite(ctx, userId, -1)
			if err2 != nil {
				return err2
			}
			return nil
		})
	default:
		return errors.New("invalid action type")
	}
}

func (uc *FavoriteUsecase) GetFavoriteList(ctx context.Context, userID uint32, tokenString string) ([]Video, error) {
	// user verification & get user_id
	token, err := common.ParseToken(uc.config.Http.TokenKey, tokenString)
	if err != nil {
		return nil, err
	}
	data, err := common.GetTokenData(token)
	if err != nil {
		return nil, err
	}
	userIDFloat64, ok := data["user_id"].(float64)
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

func (uc *FavoriteUsecase) IsFavorite(ctx context.Context, userID, videoID uint32) (bool, error) {
	// internal use; no need to verify token
	return uc.favoriteRepo.IsFavorite(ctx, userID, videoID)
}
