package data

import (
	"Atreus/app/favorite/service/internal/biz"
	"Atreus/app/favorite/service/internal/server"
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	"time"
)

// Favorite Database Model
type Favorite struct {
	ID        uint32     `gorm:"column:id;primary_key;autoIncrement"`
	VideoID   uint32     `gorm:"column:video_id"`
	UserID    uint32     `gorm:"column:user_id"`
	CreatedAt time.Time  `gorm:"column:created_at"` // new add field; for backend use only
	DeletedAt *time.Time `gorm:"column:deleted_at;index"`
}

func (Favorite) TableName() string {
	return "favorites"
}

type favoriteRepo struct {
	data        *Data
	publishRepo biz.PublishRepo
	userRepo    biz.UserRepo
	tx          biz.Transaction
	log         *log.Helper
}

func NewFavoriteRepo(
	data *Data, publishConn server.PublishConn, userConn server.UserConn, logger log.Logger) biz.FavoriteRepo {
	return &favoriteRepo{
		data:        data,
		publishRepo: NewPublishRepo(publishConn),
		userRepo:    NewUserRepo(userConn),
		log:         log.NewHelper(log.With(logger, "module", "favorite-service/repo")),
	}
}

func (r *favoriteRepo) CreateFavorite(ctx context.Context, userId, videoId uint32) error {
	// check if favorite exists
	isFavorite, err := r.IsFavorite(ctx, userId, videoId)
	if err != nil {
		return errors.New("failed to check if video is favorited")
	}
	if isFavorite {
		return errors.New("duplicate favorite(user has favoured this video)")
	}
	// create favorite
	favorite := Favorite{
		VideoID: videoId,
		UserID:  userId,
	}
	result := r.data.db.WithContext(ctx).Create(&favorite)
	if result.Error != nil {
		return fmt.Errorf("failed to create favorite: %w", result.Error)
	}
	return nil
}

// IsFavorite checks if a video is favorited by a user, avoiding duplicate favorites
func (r *favoriteRepo) IsFavorite(ctx context.Context, userId, videoId uint32) (bool, error) {
	result := r.data.db.WithContext(ctx).
		Where("user_id = ? AND video_id = ?", userId, videoId).
		First(&Favorite{})
	if result.Error == nil {
		return true, nil
	}
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return false, fmt.Errorf("failed to check if video is favorited: %w", result.Error)
}

func (r *favoriteRepo) DeleteFavorite(ctx context.Context, userId, videoId uint32) error {
	// check
	isFavorite, err := r.IsFavorite(ctx, userId, videoId)
	if err != nil {
		return errors.New("failed to check if video is favorited")
	}
	if !isFavorite {
		return errors.New("video is not favorited, failed to delete")
	}
	// delete
	result := r.data.db.WithContext(ctx).
		Where("user_id = ? AND video_id = ?", userId, videoId).
		Delete(&Favorite{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete favorite: %w", result.Error)
	}
	return nil
}

// GetFavoriteList returns a list of favorite videos(not literally the "favorite" model) of a user
func (r *favoriteRepo) GetFavoriteList(ctx context.Context, userID uint32) ([]biz.Video, error) {
	// query favorite
	var favorites []*biz.Favorite
	result := r.data.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&favorites)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get favorite list: %w", result.Error)
	}
	// convert favorite to video list
	var videoIDs []uint32
	for _, favorite := range favorites {
		videoIDs = append(videoIDs, favorite.VideoID)
	}
	videos, err := r.publishRepo.GetVideoListByVideoIds(ctx, videoIDs)
	//videos, err := nil, errors.New("not implemented")
	if err != nil {
		return nil, fmt.Errorf("failed to get video info by video ids: %w", err)
	}
	return videos, nil
}

/*func (r *favoriteRepo) CountFavoriteByVideoIDs(ctx context.Context, videoIDs []uint32) (int64, error) {
	var count int64
	result := r.data.db.WithContext(ctx).Model(&Favorite{}).Where("video_id IN (?)", videoIDs).Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to count favorite by video id: %w", result.Error)
	}
	return count, nil
}
func (r *favoriteRepo) CountFavoriteByUserID(ctx context.Context, userID uint32) (int64, error) {
	var count int64
	result := r.data.db.WithContext(ctx).Model(&Favorite{}).Where("user_id = ?", userID).Count(&count)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return 0, fmt.Errorf("failed to count favorite by user id, user id not exist: %w", result.Error)
	}
	if result.Error != nil {
		return 0, fmt.Errorf("failed to count favorite by user id: %w", result.Error)
	}
	return count, nil
}
*/
