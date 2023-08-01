package biz

import (
	"Atreus/app/comment/service/internal/conf"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewCommentUsecase)

type CommentRepo interface {
	CreateComment(context.Context, uint32, string, uint32) (*Comment, error)
	DeleteComment(context.Context, uint32, uint32, uint32) (*Comment, error)
	GetCommentList(context.Context, uint32) ([]*Comment, error)
	GetCommentNumber(context.Context, uint32) (int64, error)
}

type CommentUsecase struct {
	commentRepo CommentRepo
	config      *conf.JWT
	log         *log.Helper
}

func NewCommentUsecase(conf *conf.JWT, cr CommentRepo, logger log.Logger) *CommentUsecase {
	return &CommentUsecase{config: conf, commentRepo: cr, log: log.NewHelper(log.With(logger, "model", "usecase/comment"))}
}
