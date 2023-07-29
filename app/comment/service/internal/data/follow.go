package data

import (
	"context"
	"gorm.io/gorm"
)

// Follow Database Model
type Follow struct {
	gorm.DeletedAt
	Id         int64 `gorm:"primary_key"`
	UserId     int64 `gorm:"column:user_id;not null"`
	FollowerId int64 `gorm:"column:follower_id;not null"`
}

func (Follow) TableName() string {
	return "follows"
}

// isFollow 判断评论者是否关注视频作者
func (r *commentRepo) isFollow(ctx context.Context, userId int64, followerId int64) bool {
	result := r.data.db.WithContext(ctx).First(&Follow{}, "user_id = ? AND follower_id = ?", userId, followerId)
	return result.Error == nil
}
