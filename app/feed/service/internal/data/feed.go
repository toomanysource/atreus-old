package data

import (
	"Atreus/app/feed/service/internal/biz"
	"Atreus/app/feed/service/internal/server"
	"context"
	"strconv"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

const (
	VideoCount = 30
)

// Video DB model
type Video struct {
	Id            uint32 `gorm:"primary_key"`
	AuthorId      uint32 `gorm:"column:author_id;not null"`
	Title         string `gorm:"column:title;not null"`
	PlayURL       string `gorm:"column:play_url;not null"`
	CoverURL      string `gorm:"column:cover_url;not null"`
	FavoriteCount uint32 `gorm:"column:favorite_count;not null"`
	CommentCount  uint32 `gorm:"column:comment_count;not null"`
	// IsFavorite   s bool   // `json:"is_favorite"` // it field need to RPC. Want update.
	CreateAt string `gorm:"column:created_at;default:''"`
	gorm.DeletedAt
}

type User struct {
	*biz.User
}

func (Video) TableName() string {
	return "videos"
}

type UserRepo interface {
	GetUserInfoByUserIds(context.Context, []uint32) ([]*biz.User, error)
}
type FavoriteRepo interface {
	IsFavorite(context.Context, uint32, uint32) (bool, error)
}

type feedRepo struct {
	data     *Data
	userRepo UserRepo
	log      *log.Helper
}

func NewFeeedRepo(data *Data, userconn server.UserConn, logger log.Logger) biz.FeedRepo {
	return &feedRepo{
		data:     data,
		userRepo: NewUserRepo(userconn),
		log:      log.NewHelper(log.With(logger, "model", "data/feed")),
	}
}

// Get Feedlist no Login.
func (r *feedRepo) GetFeedList(ctx context.Context, latest_time string) (vl []*biz.Video, next_time int64, err error) {
	// check latestTime
	latestTime, err := strconv.ParseInt(latest_time, 10, 64)
	if err != nil {
		latestTime = time.Now().UnixMilli()
	}

	// // Use Cache
	// // Search if List is stored in the cache
	// feedList, err := r.data.cache.HGetAll(ctx, latest_time).Result()
	// if err != nil {
	// 	return nil, nil, fmt.Errorf("redis query error %w", err)
	// }
	// // if is stored List in cache

	var vList []Video
	// Create new feed list
	err = r.data.db.
		Where("create_at <?", latestTime).
		Order("create_at DESC").
		Limit(VideoCount).
		Find(&vList).
		Error
	if err != nil {
		return nil, 0, err
	}

	userIds := make([]uint32, 0, len(vList))
	for _, v := range vList {
		userIds = append(userIds, v.AuthorId)
	}
	users, err := r.userRepo.GetUserInfoByUserIds(ctx, userIds)
	userMap := make(map[uint32]*biz.User)
	for _, user := range users {
		userMap[user.Id] = user
	}

	for _, v := range vList {
		vl = append(vl, &biz.Video{
			Id:            v.Id,
			Author:        *userMap[v.AuthorId],
			Title:         v.Title,
			PlayURL:       v.PlayURL,
			CoverURL:      v.CoverURL,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    false,
		})
	}

	var nextTime int64
	if len(vList) > 0 {
		nextTime, err = strconv.ParseInt(vList[len(vList)-1].CreateAt, 10, 64)

		if err != nil {
			return nil, 0, err
		}
	}

	// TODO 问题：如何从Authorid链接到整个表 要怎么设计 数据结构给的是User 而不是现在使用的userid
	// 需要从vl中 获取每个video 的Author
	// 然后调用 r.userRepo.GetUserInfoByUserIDs 把内容写到vl里

	// // use Goroutine
	// // query nextTime field
	// var nextTime int64
	// if len(vl) > 0 {
	// 	err = r.data.db.
	// 		Where("create_at >?", latestTime).
	// 		Order("create_at ASC").
	// 		Limit(VideoCount).
	// 		Select("create_at").
	// 		First(&nextTime).
	// 		Error
	// }
	// if err != nil {
	// 	return nil, nil, err
	// }

	return vl, nextTime, nil
}
