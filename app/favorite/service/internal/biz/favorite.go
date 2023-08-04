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
	Id            int64
	Author        *User
	PlayUrl       string
	CoverUrl      string
	FavoriteCount int64
	CommentCount  int64
	IsFavorite    bool
	Title         string
}

// User is used to receive video info from user service;& send response
type User struct {
	Id              int64
	Name            string
	FollowCount     int64
	FollowerCount   int64
	IsFollow        bool
	Avatar          string
	BackgroundImage string
	Signature       string
	TotalFavorited  int64
	WorkCount       int64
	FavoriteCount   int64
}

// FavoriteRepo is database manipulation interface
type FavoriteRepo interface {
	CreateFavorite(context context.Context, userId, videoId uint32) error
	DeleteFavorite(context context.Context, userId, videoId uint32) error
	GetFavoriteList(context context.Context, userID uint32) ([]Video, error)           // list of user's favorite video; use slice without pointer
	CountFavoriteByVideoIDs(context context.Context, videoIDs []uint32) (int64, error) // num of favorite of a video
	CountFavoriteByUserID(context context.Context, userID uint32) (int64, error)       // num of favorite user received
}

type FavoriteUsecase struct {
	repo   FavoriteRepo
	config *conf.JWT
	log    *log.Helper
}

func NewFavoriteUsecase(repo FavoriteRepo, logger log.Logger) *FavoriteUsecase {
	return &FavoriteUsecase{repo: repo, log: log.NewHelper(logger)}
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
	userId, err := uc.getUserIDFromClaim(claims)
	if err != nil {
		return err
	}
	// biz
	switch actionType {
	case 1:
		return uc.repo.CreateFavorite(ctx, userId, videoId)
	case 2:
		return uc.repo.DeleteFavorite(ctx, userId, videoId)
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
	return uc.repo.GetFavoriteList(ctx, userID)
	//return nil, nil
}

// CountFavoriteByVideoID is for internal use; function wrap
func (uc *FavoriteUsecase) CountFavoriteByVideoID(context context.Context, videoID uint32) (int64, error) {
	return uc.repo.CountFavoriteByVideoIDs(context, []uint32{videoID})
}

// CountFavoriteByVideoIDs is for internal use
func (uc *FavoriteUsecase) CountFavoriteByVideoIDs(context context.Context, videoIDs []uint32) (int64, error) {
	count, err := uc.repo.CountFavoriteByVideoIDs(context, videoIDs)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// CountFavoriteByUserID is for internal use
func (uc *FavoriteUsecase) CountFavoriteByUserID(context context.Context, userID uint32) (int64, error) {
	count, err := uc.repo.CountFavoriteByUserID(context, userID)
	if err != nil {
		return 0, err
	}
	return count, nil
}
