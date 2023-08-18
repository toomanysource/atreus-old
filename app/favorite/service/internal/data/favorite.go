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
	ID        uint32         `gorm:"column:id;primary_key;autoIncrement"`
	VideoID   uint32         `gorm:"column:video_id"`
	UserID    uint32         `gorm:"column:user_id"`
	CreatedAt time.Time      `gorm:"column:created_at"` // new add field; for backend use only
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (Favorite) TableName() string {
	return "favorites"
}

type favoriteRepo struct {
	data        *Data
	publishRepo biz.PublishRepo
	userRepo    biz.UserRepo
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

// IsFavoriteSingle checks if a video is favorited by a user, avoiding duplicate favorites; internal use only
func (r *favoriteRepo) IsFavoriteSingle(ctx context.Context, userId, videoId uint32) (bool, error) {
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

// GetAuthorId fetch AuthorId by videoId From Publish Service
func (r *favoriteRepo) GetAuthorId(ctx context.Context, videoId uint32) (uint32, error) {
	videoList, err := r.publishRepo.GetVideoListByVideoIds(ctx, []uint32{videoId})
	if err != nil {
		return 0, errors.New("failed to fetch video author from Publish Service")
	}
	if len(videoList) == 0 {
		return 0, errors.New("video not found")
	}
	authorId := videoList[0].Author.Id
	return authorId, nil
}

// IsFavorite []videoId & userId -> []bool; exposed to biz layer
func (r *favoriteRepo) IsFavorite(ctx context.Context, userId uint32, videoIds []uint32) ([]bool, error) {
	// fetch all favorites of a user
	var favorites []Favorite
	result := r.data.db.WithContext(ctx).
		Where("user_id = ? AND video_id IN ?", userId, videoIds).
		Find(&favorites)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch favorites: %w", result.Error)
	}
	// create a map of videoId to bool
	favoriteMap := make(map[uint32]bool)
	for _, favorite := range favorites {
		favoriteMap[favorite.VideoID] = true
	}
	// create a list of bools to return
	var isFavorite []bool
	for _, videoId := range videoIds {
		if _, ok := favoriteMap[videoId]; !ok {
			isFavorite = append(isFavorite, false)
			continue
		}
		isFavorite = append(isFavorite, favoriteMap[videoId])
	}
	return isFavorite, nil
}

// CreateFavoriteTx integrates SQL database & Redis cache; exposed to biz layer
func (r *favoriteRepo) CreateFavoriteTx(ctx context.Context, userId, videoId uint32) error {
	err := r.CreateFavorite(ctx, userId, videoId)
	if err != nil {
		return err
	}
	return nil
}

// CreateFavorite IO of fav/user/video SQL database ;not exposed to biz layer
func (r *favoriteRepo) CreateFavorite(ctx context.Context, userId, videoId uint32) error {
	// check if favorite exists
	isFavorite, err := r.IsFavoriteSingle(ctx, userId, videoId)
	if err != nil {
		return errors.New("failed to check if video is favorited")
	}
	if isFavorite {
		return errors.New("duplicate favorite(user has favoured this video)")
	}
	// fetch video author id
	authorId, err := r.GetAuthorId(ctx, videoId)
	if err != nil {
		return errors.New("failed to fetch video author")
	}
	// begin transaction
	err = r.data.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// create favorite
		err = tx.Create(&Favorite{
			VideoID: videoId,
			UserID:  userId,
		}).Error
		if err != nil {
			return err
		}
		// notify other services
		err = r.userRepo.UpdateFavorited(ctx, authorId, 1)
		if err != nil {
			return fmt.Errorf("updateFavorited err: %w", err)
		}
		err = r.userRepo.UpdateFavorite(ctx, userId, 1)
		if err != nil {
			return fmt.Errorf("updateFavorite err: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to create favorite: %w", err)
	}
	return nil
}

func (r *favoriteRepo) DeleteFavoriteTx(ctx context.Context, userId, videoId uint32) error {
	err := r.DeleteFavorite(ctx, userId, videoId)
	if err != nil {
		return err
	}
	return nil
}

func (r *favoriteRepo) DeleteFavorite(ctx context.Context, userId, videoId uint32) error {
	// check if favorite exists
	isFavorite, err := r.IsFavoriteSingle(ctx, userId, videoId)
	if err != nil {
		return errors.New("failed to check if video is favorited")
	}
	if !isFavorite {
		return errors.New("video is not favorited, failed to delete")
	}
	// fetch video author id
	authorId, err := r.GetAuthorId(ctx, videoId)
	if err != nil {
		return errors.New("failed to fetch video author")
	}
	// begin transaction
	result := r.data.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// delete favorite
		err2 := tx.Where("user_id = ? AND video_id = ?", userId, videoId).Delete(&Favorite{}).Error
		if err2 != nil {
			return err2
		}
		// notify other services
		err2 = r.userRepo.UpdateFavorited(ctx, authorId, -1)
		if err2 != nil {
			return err2
		}
		err2 = r.userRepo.UpdateFavorite(ctx, userId, -1)
		if err2 != nil {
			return err2
		}
		err2 = r.publishRepo.UpdateFavoriteCount(ctx, videoId, -1)
		return nil
	})
	if result != nil {
		return errors.New("failed to delete favorite")
	}
	return nil
}

// GetFavoriteList returns a list of favorite videos(not literally the "favorite" model) of a user
func (r *favoriteRepo) GetFavoriteList(ctx context.Context, userID uint32) ([]biz.Video, error) {
	// query favorite
	var favorites []Favorite
	result := r.data.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&favorites)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get favorite list: %w", result.Error)
	}
	if len(favorites) == 0 {
		return nil, nil
	}
	// convert favorite to video list
	var videoIDs []uint32
	for _, favorite := range favorites {
		videoIDs = append(videoIDs, favorite.VideoID)
	}
	videos, err := r.publishRepo.GetVideoListByVideoIds(ctx, videoIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get video info by video ids: %w", err)
	}
	for _, video := range videos {
		video.IsFavorite = true
	}
	return videos, nil
}
