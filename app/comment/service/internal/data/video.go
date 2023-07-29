package data

import (
	"context"
	"gorm.io/gorm"
)

// Video Database Model
type Video struct {
	gorm.Model
	Id            int64  `gorm:"column:id;primary_key"`
	AuthorId      int64  `gorm:"column:author_id;not null"`
	Title         string `gorm:"column:title;not null"`
	PlayUrl       string `gorm:"column:play_url;not null"`
	FavoriteCount int64  `gorm:"column:favorite_count;not null;default:0"`
	CommentCount  int64  `gorm:"column:comment_count;not null;default:0"`
	CreatedAt     string `gorm:"column:created_at;not null;default:''"`
}

func (Video) TableName() string {
	return "videos"
}

// getVideoAuthor 获取视频作者信息
func (r *commentRepo) getVideoAuthor(ctx context.Context, videoId int64) (*User, error) {
	var video = &Video{}
	result := r.data.db.First(video, videoId)
	if err := result.Error; err != nil {
		return nil, err
	}
	return r.getUser(ctx, video.AuthorId)
}
