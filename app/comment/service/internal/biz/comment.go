package biz

import (
	"Atreus/app/comment/service/internal/conf"
	"Atreus/pkg/common"
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
)

type Comment struct {
	Id         uint32
	User       *User
	Content    string
	CreateDate string
}

type User struct {
	Id              uint32
	Name            string
	Avatar          string
	BackgroundImage string
	Signature       string
	IsFollow        bool
	FollowCount     uint32
	FollowerCount   uint32
	TotalFavorited  uint32
	WorkCount       uint32
	FavoriteCount   uint32
}

type CommentRepo interface {
	CreateComment(context.Context, uint32, string, uint32) (*Comment, error)
	DeleteComment(context.Context, uint32, uint32, uint32) (*Comment, error)
	GetCommentList(context.Context, uint32, uint32) ([]*Comment, error)
}

type CommentUsecase struct {
	commentRepo CommentRepo
	config      *conf.JWT
	log         *log.Helper
}

func NewCommentUsecase(conf *conf.JWT, cr CommentRepo, logger log.Logger) *CommentUsecase {
	return &CommentUsecase{config: conf, commentRepo: cr, log: log.NewHelper(log.With(logger, "model", "usecase/comment"))}
}

func (uc *CommentUsecase) GetCommentList(
	ctx context.Context, tokenString string, videoId uint32) ([]*Comment, error) {
	// 未登录状态
	if tokenString != "" {
		token, err := common.ParseToken(uc.config.Http.TokenKey, tokenString)
		if err != nil {
			return nil, err
		}
		data, err := common.GetTokenData(token)
		if err != nil {
			return nil, err
		}
		userId := uint32(data["user_id"].(float64))
		return uc.commentRepo.GetCommentList(ctx, userId, videoId)
	}
	return uc.commentRepo.GetCommentList(ctx, 0, videoId)
}

func (uc *CommentUsecase) CommentAction(
	ctx context.Context, videoId, commentId uint32,
	actionType uint32, commentText string, tokenString string) (*Comment, error) {
	token, err := common.ParseToken(uc.config.Http.TokenKey, tokenString)
	if err != nil {
		return nil, err
	}
	data, err := common.GetTokenData(token)
	if err != nil {
		return nil, err
	}
	userId := uint32(data["user_id"].(float64))
	switch actionType {
	case 1:
		return uc.commentRepo.CreateComment(ctx, videoId, commentText, userId)
	case 2:
		return uc.commentRepo.DeleteComment(ctx, videoId, commentId, userId)
	default:
		return nil, errors.New("the value of action_type is not in the specified range")
	}
}
