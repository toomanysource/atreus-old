package data

import (
	"Atreus/app/feed/service/internal/biz"
	"Atreus/app/feed/service/internal/server"
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

var (
	VideoCount uint32 = 30
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
	CreatedAt string `gorm:"column:created_at;default:''"`
	gorm.DeletedAt
}

func (Video) TableName() string {
	return "videos"
}

//	type UserRepo interface {
//		GetUserInfoByUserIds(context.Context, []uint32) ([]*biz.User, error)
//	}
type PublishRepo interface {
	GetVideoList(context.Context, string, uint32, uint32) (int64, []*Video, error)
}

type feedRepo struct {
	data        *Data
	publishRepo biz.PublishRepo
	log         *log.Helper
}

func NewFeeedRepo(data *Data, publishconn server.PublishConn, logger log.Logger) biz.FeedRepo {
	return &feedRepo{
		data:        data,
		publishRepo: NewPublishRepo(publishconn),
		log:         log.NewHelper(log.With(logger, "model", "data/feed")),
	}
}

func (r *feedRepo) GetFeedListById(ctx context.Context, latest_time string, user_id uint32) (next_time int64, vl []biz.Video, err error) {
	if latest_time == "" {
		latest_time = time.Now().String()
	}
	nextTime, vl, err := r.publishRepo.GetVideoList(ctx, latest_time, user_id, VideoCount)
	if err != nil {
		return 0, nil, err
	}
	return nextTime, vl, nil
	// return vl, nextTime, nil
}

// Get Feedlist no Login.
func (r *feedRepo) GetFeedList(ctx context.Context, latest_time string) (next_time int64, vl []biz.Video, err error) {
	// check latestTime
	// latestTime, err := strconv.ParseInt(latest_time, 10, 64)
	// if err != nil {
	// 	latestTime = time.Now().UnixMilli()
	// }
	if latest_time == "" {
		latest_time = time.Now().String()
	}
	nextTime, vl, err := r.publishRepo.GetVideoList(ctx, latest_time, 0, VideoCount)
	if err != nil {
		return 0, nil, err
	}
	return nextTime, vl, nil
	// var vList []Video
	// // Create new feed list
	// err = r.data.db.
	// 	Where("created_at <?", latestTime).
	// 	Order("created_at DESC").
	// 	Limit(VideoCount).
	// 	Find(&vList).
	// 	Error
	// if err != nil {
	// 	return nil, 0, err
	// }

	// userIds := make([]uint32, 0, len(vList))
	// for _, v := range vList {
	// 	userIds = append(userIds, v.AuthorId)
	// }
	// users, err := r.userRepo.GetUserInfoByUserIds(ctx, userIds)
	// if err != nil {
	// 	return nil, 0, err
	// }
	// userMap := make(map[uint32]*biz.User)
	// for _, user := range users {
	// 	userMap[user.Id] = user
	// }

	// for _, v := range vList {
	// 	vl = append(vl, &biz.Video{
	// 		Id:            v.Id,
	// 		Author:        *userMap[v.AuthorId],
	// 		Title:         v.Title,
	// 		PlayURL:       v.PlayURL,
	// 		CoverURL:      v.CoverURL,
	// 		FavoriteCount: v.FavoriteCount,
	// 		CommentCount:  v.CommentCount,
	// 		IsFavorite:    false,
	// 	})
	// }

	// var nextTime int64
	// if len(vList) > 0 {
	// 	nextTime, err = strconv.ParseInt(vList[len(vList)-1].CreatedAt, 10, 64)

	// 	if err != nil {
	// 		return nil, 0, err
	// 	}
	// }
	// // // use Goroutine

	// return vl, nextTime, nil
}
