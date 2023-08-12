package biz

import (
	"Atreus/app/favorite/service/internal/conf"
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-jwt/jwt/v4"
	"strconv"
)

// Favorite is corresponding to the favorite table in database
type Favorite struct {
	//ID      int -> use <Composite Primary Key> instead
	VideoID uint32
	UserID  uint32
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

func NewFavoriteUsecase(repo FavoriteRepo, logger log.Logger) *FavoriteUsecase {
	return &FavoriteUsecase{favoriteRepo: repo, log: log.NewHelper(logger)}
}

// parseToken verify token & return claims
func (uc *FavoriteUsecase) parseToken(tokenKey, tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.MapClaims{},
		func(token *jwt.Token) (interface{}, error) { return []byte(tokenKey), nil })
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

// getUserIDFromToken : extract user_id from jwt.MapClaims
func (uc *FavoriteUsecase) getUserIDFromClaim(claims *jwt.MapClaims) (uint32, error) {
	id, err := strconv.ParseUint((*claims)["user_id"].(string), 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(id), nil
}

// FavoriteAction is for http api use; create & delete integrated
func (uc *FavoriteUsecase) FavoriteAction(ctx context.Context, videoId, actionType uint32, tokenString string) error {
	// user verification & get user_id
	claims, err := uc.parseToken(uc.config.Http.TokenKey, tokenString)
	if err != nil {
		return err
	}
	userId, err := uc.getUserIDFromClaim(claims) // to update user's favorite
	if err != nil {
		return err
	}
	// biz
	videoInfo, err := uc.publishRepo.GetVideoListByVideoIds(ctx, []uint32{videoId})
	if err != nil {
		return err
	}
	if len(videoInfo) == 0 {
		return errors.New("video not found")
	}
	authorId := videoInfo[0].Author.Id // to update author's favorited
	//err =
	//return err
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
	claims, err := uc.parseToken(uc.config.Http.TokenKey, tokenString)
	if err != nil {
		return nil, err
	}
	userId, err := uc.getUserIDFromClaim(claims)
	if err != nil {
		return nil, err
	}
	if userId != userID {
		return nil, errors.New("jwt and request user_id not match")
	}
	// biz
	return uc.favoriteRepo.GetFavoriteList(ctx, userID)
	//return nil, nil
}

func (uc *FavoriteUsecase) IsFavorite(ctx context.Context, userID, videoID uint32) (bool, error) {
	// internal use; no need to verify token
	return uc.favoriteRepo.IsFavorite(ctx, userID, videoID)
}

// CountFavoriteByVideoID is for internal use; function wrap
//func (uc *FavoriteUsecase) CountFavoriteByVideoID(context context.Context, videoID uint32) (int64, error) {
//	return uc.favoriteRepo.CountFavoriteByVideoIDs(context, []uint32{videoID})
//}
//
//// CountFavoriteByVideoIDs is for internal use
//func (uc *FavoriteUsecase) CountFavoriteByVideoIDs(context context.Context, videoIDs []uint32) (int64, error) {
//	count, err := uc.favoriteRepo.CountFavoriteByVideoIDs(context, videoIDs)
//	if err != nil {
//		return 0, err
//	}
//	return count, nil
//}
//
//// CountFavoriteByUserID is for internal use
//func (uc *FavoriteUsecase) CountFavoriteByUserID(context context.Context, userID uint32) (int64, error) {
//	count, err := uc.favoriteRepo.CountFavoriteByUserID(context, userID)
//	if err != nil {
//		return 0, err
//	}
//	return count, nil
//}
