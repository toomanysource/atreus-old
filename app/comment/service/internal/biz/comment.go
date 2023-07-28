package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

// Comment Data Model
type Comment struct {
	Id         int64
	User       *User
	Content    string
	CreateDate string
}

type User struct {
	Id              int64
	Name            string
	Avatar          string
	BackgroundImage string
	Signature       string
	IsFollow        bool
	FollowCount     int64
	FollowerCount   int64
	TotalFavorited  int64
	WorkCount       int64
	FavoriteCount   int64
}

type CommentRepo interface {
	PublishComment(context.Context, int64, string, string) (*Comment, error)
	DeleteComment(context.Context, int64, int64, string) (*Comment, error)
	GetCommentList(context.Context, int64, string) ([]*Comment, error)
}

type CommentUsecase struct {
	repo CommentRepo
	log  *log.Helper
}

func NewCommentUsecase(repo CommentRepo, logger log.Logger) *CommentUsecase {
	return &CommentUsecase{repo: repo, log: log.NewHelper(log.With(logger, "model", "usecase/comment"))}
}

func (uc *CommentUsecase) GetCommentList(
	ctx context.Context, token string, videoId int64) ([]*Comment, error) {
	return uc.repo.GetCommentList(ctx, videoId, token)
}

func (uc *CommentUsecase) CommentAction(
	ctx context.Context, videoId, commentId int64,
	actionType int32, commentText, token string) (c *Comment, err error) {
	if actionType == 1 {
		return uc.repo.PublishComment(ctx, videoId, commentText, token)
	} else if actionType == 2 {
		return uc.repo.DeleteComment(ctx, videoId, commentId, token)
	}
	return nil, err
}
