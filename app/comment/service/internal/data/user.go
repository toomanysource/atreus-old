package data

import (
	"context"
	"gorm.io/gorm"
)

// User Database Model
type User struct {
	gorm.Model
	Id                 int64  `gorm:"column:id;primary_key"`
	Name               string `gorm:"column:name;not null"`
	Password           string `gorm:"column:password;not null"`
	AvatarUrl          string `gorm:"column:avatar_url;not null;default:''"`
	BackgroundImageUrl string `gorm:"column:background_image_url;not null;default:''"`
	Signature          string `gorm:"column:signature;not null;default:''"`
	FollowCount        int64  `gorm:"column:follow_count;not null;default:0"`
	FollowerCount      int64  `gorm:"column:follower_count;not null;default:0"`
	TotalFavorited     int64  `gorm:"column:total_favorited;not null;default:0"`
	WorkCount          int64  `gorm:"column:work_count;not null;default:0"`
	FavoriteCount      int64  `gorm:"column:favorite_count;not null;default:0"`
}

func (User) TableName() string {
	return "users"
}

// getUser 用户表获取用户信息
func (r *commentRepo) getUser(ctx context.Context, userId int64) (*User, error) {
	var user = &User{}
	result := r.data.db.WithContext(ctx).First(user, userId)
	if err := result.Error; err != nil {
		return nil, err
	}
	return user, nil
}
