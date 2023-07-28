package data

import (
	"context"
	"errors"
	"time"

	"Atreus/app/comment/service/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type commentRepo struct {
	data *Data
	log  *log.Helper
}

func NewCommentRepo(data *Data, logger log.Logger) biz.CommentRepo {
	return &commentRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "model", "comment-service/repo")),
	}
}

// Comment Database Model
type Comment struct {
	Id         int64  `gorm:"primary_key"`
	UserId     int64  `gorm:"column:user_id;not null"`
	VideoId    int64  `gorm:"column:video_id;not null"`
	Content    string `gorm:"column:content;not null"`
	CreateDate string `gorm:"column:create_date;default:''"`
}

type User struct {
	id                 int64
	name               string
	password           string
	avatarUrl          string
	backgroundImageUrl string
	signature          string
	followCount        int64
	followerCount      int64
	totalFavorited     int64
	workCount          int64
	favoriteCount      int64
	createDate         string
	update_date        string
	delete_date        string
}

type Video struct {
	id            int64
	authorId      int64
	title         string
	playUrl       string
	favoriteCount int64
	commentCount  int64
	createDate    string
}

type Follow struct {
	id         int64
	userId     int64
	followerId int64
}

func (r *commentRepo) PublishComment(
	ctx context.Context, videoId, userId int64, commentText string) (c *biz.Comment, err error) {
	return r.createComment(ctx, videoId, userId, commentText)
}

func (r *commentRepo) DeleteComment(
	ctx context.Context, videoId, commentId int64, token string) (c *biz.Comment, err error) {
	return r.deleteComment(ctx, videoId, commentId, token)
}

func (r *commentRepo) GetCommentList(ctx context.Context, videoId int64, token string) (cl []*biz.Comment, err error) {
	return nil, nil
}

func (r *commentRepo) deleteComment(
	ctx context.Context, videoId, commentId int64, token string) (c *biz.Comment, err error) {
	comment := &Comment{}
	result := r.data.db.WithContext(ctx).First(comment, commentId)
	if err = result.Error; err != nil {
		return nil, err
	}
	if comment.UserId != token.user_id {
		return nil, errors.New("mismatch between commenter id and user id.")
	}
	if comment.VideoId != videoId {
		return nil, errors.New("comment video id doesn't match current video id.")
	}
	result = r.data.db.WithContext(ctx).Delete(&Comment{}, commentId)
	if err = result.Error; err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *commentRepo) createComment(
	ctx context.Context, videoId, userId int64, commentText string) (c *biz.Comment, err error) {
	comment := &Comment{
		UserId:     userId,
		VideoId:    videoId,
		Content:    commentText,
		CreateDate: time.Now().Format("01-02"),
	}
	result := r.data.db.WithContext(ctx).Create(comment)
	if err = result.Error; err != nil {
		return nil, err
	}
	user, err := r.getCommentUser(ctx, userId, videoId)
	if err != nil {
		return nil, err
	}
	return &biz.Comment{
		Id:         comment.Id,
		User:       user,
		Content:    commentText,
		CreateDate: comment.CreateDate,
	}, nil
}

func (r *commentRepo) getCommentUser(ctx context.Context, userId, videoId int64) (u *biz.User, err error) {
	user, err := r.getUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	author, err := r.getVideoAuthor(ctx, videoId)
	return &biz.User{
		Id:              user.id,
		Name:            user.name,
		Avatar:          user.avatarUrl,
		BackgroundImage: user.backgroundImageUrl,
		Signature:       user.signature,
		FollowCount:     user.followCount,
		FollowerCount:   user.followerCount,
		TotalFavorited:  user.totalFavorited,
		WorkCount:       user.workCount,
		FavoriteCount:   user.favoriteCount,
		IsFollow:        r.isFollow(ctx, author.id, userId),
	}, nil
}

func (r *commentRepo) isFollow(ctx context.Context, userId int64, followerId int64) bool {
	result := r.data.db.WithContext(ctx).First(&Follow{}, "user_id = ? AND follower_id = ?", userId, followerId)
	if result.Error != nil {
		return false
	}
	return true
}

func (r *commentRepo) getUser(ctx context.Context, userId int64) (u *User, err error) {
	var user = &User{}
	result := r.data.db.WithContext(ctx).First(user, userId)
	if err = result.Error; err != nil {
		return nil, err
	}
	return user, nil
}

// getVideoAuthor
func (r *commentRepo) getVideoAuthor(ctx context.Context, videoId int64) (u *User, err error) {
	var video = &Video{}
	result := r.data.db.First(video, videoId)
	if err = result.Error; err != nil {
		return nil, err
	}
	return r.getUser(ctx, video.authorId)
}
