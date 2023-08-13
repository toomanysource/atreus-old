package data

import (
	"Atreus/app/publish/service/internal/biz"
	"Atreus/app/publish/service/internal/server"
	"Atreus/pkg/ffmpegX"
	"bytes"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/minio/minio-go/v7"
	"io"
	"os"
	"strconv"
	"time"
)

// Video Database Model
type Video struct {
	Id            uint32 `gorm:"column:id;primary_key;auto_increment"`
	AuthorID      uint32 `gorm:"column:author_id;not null"`
	Title         string `gorm:"column:title;not null;size:255"`
	PlayURL       string `gorm:"column:play_url;not null"`
	CoverURL      string `gorm:"column:cover_url;not null"`
	FavoriteCount uint32 `gorm:"column:favorite_count;not null;default:0"`
	CommentCount  uint32 `gorm:"column:comment_count;not null;default:0"`
	CreateAt      int64  `gorm:"column:create_at"`
}

type UserRepo interface {
	GetUserInfos(context.Context, []uint32) ([]*biz.User, error)
}
type FavoriteRepo interface {
	IsFavorite(context.Context, uint32, []uint32) ([]bool, error)
}

type publishRepo struct {
	data         *Data
	favoriteRepo FavoriteRepo
	userRepo     UserRepo
	log          *log.Helper
}

func NewPublishRepo(
	data *Data, userConn server.UserConn, favoriteConn server.FavoriteConn, logger log.Logger) biz.PublishRepo {
	return &publishRepo{
		data:         data,
		favoriteRepo: NewFavoriteRepo(favoriteConn),
		userRepo:     NewUserRepo(userConn),
		log:          log.NewHelper(logger),
	}
}

// UploadVideo 上传视频
func (r *publishRepo) UploadVideo(ctx context.Context, fileBytes []byte, userId uint32, title string) error {
	reader := bytes.NewReader(fileBytes)
	// 生成封面
	coverReader, err := r.GenerateCoverImage(fileBytes)
	cover := coverReader.(*bytes.Reader)
	if err != nil {
		return fmt.Errorf("generate cover image error: %w", err)
	}
	// 上传封面
	err = r.data.oss.UploadSizeFile(
		ctx, "oss", title+".png", cover, cover.Size(), minio.PutObjectOptions{
			ContentType: "image/png",
		})
	// 上传视频
	err = r.data.oss.UploadSizeFile(ctx, "oss", title+".mp4", reader, reader.Size(), minio.PutObjectOptions{
		ContentType: "video/mp4",
	})
	if err != nil {
		return fmt.Errorf("upload video error: %w", err)
	}
	// 获取视频和封面的url
	playURL, coverURL, err := r.GetRemoteVideoInfo(ctx, title)
	if err != nil {
		return fmt.Errorf("get remote video info error: %w", err)
	}
	v := &Video{
		AuthorID:      userId,
		Title:         title,
		PlayURL:       playURL,
		CoverURL:      coverURL,
		FavoriteCount: 0,
		CommentCount:  0,
		CreateAt:      time.Now().Unix(),
	}
	if err := r.data.db.Create(v).Error; err != nil {
		return fmt.Errorf("create video error: %w", err)
	}
	return nil
}

// GetRemoteVideoInfo 获取远程视频信息
func (r *publishRepo) GetRemoteVideoInfo(ctx context.Context, title string) (playURL, coverURL string, err error) {
	url, err := r.data.oss.GetFileURL(
		ctx, "oss", "video/"+title+".mp4", time.Hour*24*7)
	if err != nil {
		return "", "", fmt.Errorf("get video url error: %w", err)
	}
	playURL = url.String()
	url, err = r.data.oss.GetFileURL(
		ctx, "oss", "image/"+title+".png", time.Hour*24*7)
	if err != nil {
		return "", "", fmt.Errorf("get image url error: %w", err)
	}
	coverURL = url.String()
	return
}

// GenerateCoverImage 生成封面
func (r *publishRepo) GenerateCoverImage(fileBytes []byte) (io.Reader, error) {
	// 创建临时文件
	tempFile, err := os.CreateTemp("", "tempFile-*")
	if err != nil {
		return nil, fmt.Errorf("create temp file error: %w", err)
	}
	defer os.Remove(tempFile.Name())
	if _, err = tempFile.Write(fileBytes); err != nil {
		return nil, fmt.Errorf("write temp file error: %w", err)
	}
	// 调用ffmpeg 生成封面
	return ffmpegX.ReadFrameAsImage(tempFile.Name(), 60)
}

func (r *publishRepo) FindVideoListByUserId(ctx context.Context, userId uint32) ([]*biz.Video, error) {
	var videoList []*Video
	var vl []*biz.Video
	err := r.data.db.WithContext(ctx).Where("author_id = ?", userId).Find(&videoList).Error
	if err != nil {
		return nil, err
	}
	users, err := r.userRepo.GetUserInfos(ctx, []uint32{userId})
	if err != nil {
		return nil, err
	}
	for _, video := range videoList {
		vl = append(vl, &biz.Video{
			ID:            video.Id,
			Author:        users[0],
			PlayURL:       video.PlayURL,
			CoverURL:      video.CoverURL,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    false,
			Title:         video.Title,
		})
	}
	return vl, err
}

func (r *publishRepo) FindVideoListByIDs(ctx context.Context, ids []uint32) ([]*biz.Video, error) {
	var videoList []*Video
	err := r.data.db.WithContext(ctx).Find(&videoList, ids).Error
	if err != nil {
		return nil, err
	}
	return r.GetUsers(ctx, videoList)
}

func (r *publishRepo) UpdateFavoriteCount(ctx context.Context, videoId uint32, favoriteChange int32) error {
	var video Video
	err := r.data.db.WithContext(ctx).Where("id = ?", videoId).First(&video).Error
	if err != nil {
		return err
	}
	newCount := calculateValidUint32(video.FavoriteCount, favoriteChange)
	err = r.data.db.WithContext(ctx).Model(&Video{}).Where("id = ?", videoId).Update("favorite_count", newCount).Error
	if err != nil {
		return err
	}
	return err
}

func (r *publishRepo) UpdateCommentCount(ctx context.Context, videoId uint32, commentChange int32) error {
	var video Video
	err := r.data.db.WithContext(ctx).Where("id = ?", videoId).First(&video).Error
	if err != nil {
		return err
	}
	newCount := calculateValidUint32(video.FavoriteCount, commentChange)
	err = r.data.db.WithContext(ctx).Model(&Video{}).Where("id = ?", videoId).
		Update("comment_count", newCount).Error
	if err != nil {
		return err
	}
	return err
}

func (r *publishRepo) FindVideoListByTime(
	ctx context.Context, latestTime string, userId uint32, number uint32) (int64, []*biz.Video, error) {
	var videoList []*Video
	times, err := strconv.Atoi(latestTime)
	if err != nil {
		return 0, nil, err
	}
	err = r.data.db.WithContext(ctx).Where("created_at < ?", uint64(times)).
		Order("created_at desc").Limit(int(number)).Find(&videoList).Error
	if err != nil {
		return 0, nil, err
	}
	nextTime := videoList[len(videoList)-1].CreateAt
	videoIds := make([]uint32, 0, len(videoList))
	for _, video := range videoList {
		videoIds = append(videoIds, video.Id)
	}
	isFavoriteList, err := r.favoriteRepo.IsFavorite(ctx, userId, videoIds)
	if err != nil {
		return 0, nil, err
	}
	vl, err := r.GetUsers(ctx, videoList)
	if err != nil {
		return 0, nil, err
	}
	for i, video := range vl {
		video.IsFavorite = isFavoriteList[i]
	}
	return nextTime, vl, err
}

func (r *publishRepo) GetUsers(ctx context.Context, videoList []*Video) (vl []*biz.Video, err error) {
	userIds := make([]uint32, 0, len(videoList))
	for _, video := range videoList {
		userIds = append(userIds, video.AuthorID)
	}
	userMap := make(map[uint32]*biz.User)
	users, err := r.userRepo.GetUserInfos(ctx, userIds)
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		userMap[user.ID] = user
	}
	for _, video := range videoList {
		vl = append(vl, &biz.Video{
			ID:            video.Id,
			Author:        userMap[video.AuthorID],
			PlayURL:       video.PlayURL,
			CoverURL:      video.CoverURL,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    true,
			Title:         video.Title,
		})
	}
	return
}

func calculateValidUint32(src uint32, mod int32) uint32 {
	if mod < 0 {
		mod = -mod
		if src < uint32(mod) {
			return 0
		}
		return src - uint32(mod)
	}
	return src + uint32(mod)
}
