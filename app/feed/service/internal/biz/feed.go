package biz

import (
	"Atreus/app/feed/service/internal/conf"
	"Atreus/pkg/common"
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

// var (
// 	// ErrUserNotFound is user not found.
// 	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
// )

// Feed is a feed model.
type Feed struct {
	NextTime  *int64  `json:"next_time"`
	VideoList []Video `json:"video_list"`
}

type Video struct {
	Id            uint32 `json:"id"`
	Author        User   `json:"author"`
	CommentCount  uint32 `json:"comment_count"`
	CoverURL      string `json:"cover_url"`
	FavoriteCount uint32 `json:"favorite_count"`
	IsFavorite    bool   `json:"is_favorite"`
	PlayURL       string `json:"play_url"`
	Title         string `json:"title"`
}

type User struct {
	Id              uint32 `json:"id"`
	Name            string `json:"name"`
	Avatar          string `json:"avatar"`
	BackgroundImage string `json:"background_image"`
	FavoriteCount   uint32 `json:"favorite_count"`
	FollowCount     uint32 `json:"follow_count"`
	FollowerCount   uint32 `json:"follower_count"`
	IsFollow        bool   `json:"is_follow"`
	Signature       string `json:"signature"`
	TotalFavorited  uint32 `json:"total_favorited"`
	WorkCount       uint32 `json:"work_count"`
}

// FeedRepo is a feed repo.
type FeedRepo interface {
	// FeedList(context.Context, string, uint32) ([]*Video, int64, error)
	GetFeedList(context.Context, string) (int64, []Video, error)
	GetFeedListById(context.Context, string, uint32) (int64, []Video, error)
}
type PublishRepo interface {
	GetVideoList(ctx context.Context, latest_time string, user_id uint32, number uint32) (int64, []Video, error)
}

// FeedUsecase is a feed usecase.
type FeedUsecase struct {
	repo FeedRepo
	// cache  FeedCache
	config *conf.JWT
	log    *log.Helper
}

// NewFeedUsecase new a feed usecase.
func NewFeedUsecase(repo FeedRepo, conf *conf.JWT, logger log.Logger) *FeedUsecase {
	return &FeedUsecase{repo: repo, config: conf, log: log.NewHelper(log.With(logger, "model", "usecase/feed"))}
}

// FeedList .
func (uc *FeedUsecase) FeedList(ctx context.Context, latest_time string, tokenString string) (int64, []Video, error) {
	if tokenString == "" {
		return uc.repo.GetFeedList(ctx, latest_time)
	}
	token, err := common.ParseToken(uc.config.Http.TokenKey, tokenString)
	if err != nil {
		return 0, nil, err
	}
	data, err := common.GetTokenData(token)
	if err != nil {
		return 0, nil, err
	}
	userId := uint32(data["user_id"].(float64))
	return uc.repo.GetFeedListById(ctx, latest_time, userId)
}

// // GetFeedList .
// func (uc *FeedUsecase) GetFeedList(ctx context.Context, latest_time string) ([]*Video, int64, error) {
// 	return uc.repo.GetFeedList(ctx, latest_time)
// }
