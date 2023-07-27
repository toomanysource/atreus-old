package data

import (
	"context"

	"Atreus/app/comment/service/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type commentRepo struct {
	data *Data
	log  *log.Helper
}

// NewCommentRepo .
func NewCommentRepo(data *Data, logger log.Logger) biz.CommentRepo {
	return &commentRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *commentRepo) Save(ctx context.Context, g *biz.Comment) (*biz.Comment, error) {
	return g, nil
}

func (r *commentRepo) Update(ctx context.Context, g *biz.Comment) (*biz.Comment, error) {
	return g, nil
}

func (r *commentRepo) FindByID(context.Context, int64) (*biz.Comment, error) {
	return nil, nil
}

func (r *commentRepo) ListByHello(context.Context, string) ([]*biz.Comment, error) {
	return nil, nil
}

func (r *commentRepo) ListAll(context.Context) ([]*biz.Comment, error) {
	return nil, nil
}
