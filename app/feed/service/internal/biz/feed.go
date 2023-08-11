package biz

import (
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
	// GetFeedList(context.Context, string) (*Feed, error)
	GetFeedList(context.Context, string) ([]*Video, int64, error)

	// GetFeedListById(context.Context, string, uint32) ([]*Video, *int64, error)
}

// FeedUsecase is a feed usecase.
type FeedUsecase struct {
	repo FeedRepo
	// cache  FeedCache
	// config *conf.JWT
	log *log.Helper
}

// NewFeedUsecase new a feed usecase.
func NewFeedUsecase(repo FeedRepo, logger log.Logger) *FeedUsecase {
	return &FeedUsecase{repo: repo, log: log.NewHelper(logger)}
}

// GetFeedList .

// 只需写逻辑 例如判断通过了执行什么 剩下的执行和操作代码在data层编写
func (uc *FeedUsecase) GetFeedList(ctx context.Context, latest_time string) ([]*Video, int64, error) {
	// latestTime, err := strconv.ParseInt(latest_time, 10, 64)
	// if err != nil {
	// 	return nil, nil, err
	// }

	return uc.repo.GetFeedList(ctx, latest_time)
}
