package data

import (
	"context"
	"errors"
	"gorm.io/gorm"
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
	Id       int64  `gorm:"primary_key"`
	UserId   int64  `gorm:"column:user_id;not null"`
	VideoId  int64  `gorm:"column:video_id;not null"`
	Content  string `gorm:"column:content;not null"`
	CreateAt string `gorm:"column:created_at;default:''"`
	gorm.DeletedAt
}

func (Comment) TableName() string {
	return "comments"
}

// DeleteComment 删除评论
func (r *commentRepo) DeleteComment(
	ctx context.Context, videoId, commentId int64, userId int64) (c *biz.Comment, err error) {
	comment := &Comment{}
	result := r.data.db.WithContext(ctx).First(comment, commentId)
	if err = result.Error; err != nil {
		return nil, err
	}

	// 判断当前用户是否为评论用户
	if comment.UserId != userId {
		return nil, errors.New("mismatch between commenter id and user id")
	}
	// 判断视频id是否为当前视频id
	if comment.VideoId != videoId {
		return nil, errors.New("comment video id doesn't match current video id")
	}

	result = r.data.db.WithContext(ctx).Delete(&Comment{}, commentId)
	if err = result.Error; err != nil {
		return nil, err
	}
	return nil, nil
}

// CreateComment 创建评论
func (r *commentRepo) CreateComment(
	ctx context.Context, videoId int64, commentText string, userId int64) (*biz.Comment, error) {
	comment := &Comment{
		UserId:   userId,
		VideoId:  videoId,
		Content:  commentText,
		CreateAt: time.Now().Format("01-02"),
	}

	result := r.data.db.WithContext(ctx).Create(comment)
	if err := result.Error; err != nil {
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
		CreateDate: comment.CreateAt,
	}, nil
}

// GetCommentList 获取评论列表
func (r *commentRepo) GetCommentList(ctx context.Context, videoId int64, userId int64) (cl []*biz.Comment, err error) {
	var commentList []*Comment
	r.data.db.WithContext(ctx).Where("video_id = ?", videoId).Find(commentList)
	for _, comment := range commentList {
		user, err := r.getCommentUser(ctx, comment.UserId, videoId)
		if err != nil {
			return nil, err
		}
		cl = append(cl, &biz.Comment{
			Id:         comment.Id,
			User:       user,
			Content:    comment.Content,
			CreateDate: comment.CreateAt,
		})
	}
	return
}

// getCommentUser 获取评论者信息
func (r *commentRepo) getCommentUser(ctx context.Context, userId, videoId int64) (*biz.User, error) {
	user, err := r.getUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	author, err := r.getVideoAuthor(ctx, videoId)
	if err != nil {
		return nil, err
	}

	return &biz.User{
		Id:              user.Id,
		Name:            user.Name,
		Avatar:          user.AvatarUrl,
		BackgroundImage: user.BackgroundImageUrl,
		Signature:       user.Signature,
		FollowCount:     user.FollowCount,
		FollowerCount:   user.FollowerCount,
		TotalFavorited:  user.TotalFavorited,
		WorkCount:       user.WorkCount,
		FavoriteCount:   user.FavoriteCount,
		IsFollow:        r.isFollow(ctx, author.Id, userId),
	}, nil
}
