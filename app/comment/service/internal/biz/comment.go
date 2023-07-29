package biz

import (
	"context"
	"errors"

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
	CreateComment(context.Context, int64, string, int64) (*Comment, error)
	DeleteComment(context.Context, int64, int64, int64) (*Comment, error)
	GetCommentList(context.Context, int64, int64) ([]*Comment, error)
}

type CommentUsecase struct {
	repo CommentRepo
	log  *log.Helper
}

func NewCommentUsecase(repo CommentRepo, logger log.Logger) *CommentUsecase {
	return &CommentUsecase{repo: repo, log: log.NewHelper(log.With(logger, "model", "usecase/comment"))}
}

func (uc *CommentUsecase) GetCommentList(
	ctx context.Context, tokenUserId int64, videoId int64) ([]*Comment, error) {
	return uc.repo.GetCommentList(ctx, videoId, tokenUserId)
}

func (uc *CommentUsecase) CommentAction(
	ctx context.Context, videoId, commentId int64,
	actionType int32, commentText string, tokenUserId int64) (*Comment, error) {
	if actionType == 1 {
		return uc.repo.CreateComment(ctx, videoId, commentText, tokenUserId)
	} else if actionType == 2 {
		return uc.repo.DeleteComment(ctx, videoId, commentId, tokenUserId)
	}
	return nil, errors.New("the value of action_type is not in the specified range")
}
