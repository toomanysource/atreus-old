package data

import (
	"context"

	"Atreus/app/feed/service/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type feedRepo struct {
	data *Data
	log  *log.Helper
}

// NewFeedRepo .
func NewFeedRepo(data *Data, logger log.Logger) biz.FeedRepo {
	return &feedRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *feedRepo) Save(ctx context.Context, g *biz.Feed) (*biz.Feed, error) {
	return g, nil
}

func (r *feedRepo) Update(ctx context.Context, g *biz.Feed) (*biz.Feed, error) {
	return g, nil
}

func (r *feedRepo) FindByID(context.Context, int64) (*biz.Feed, error) {
	return nil, nil
}

func (r *feedRepo) ListByHello(context.Context, string) ([]*biz.Feed, error) {
	return nil, nil
}

func (r *feedRepo) ListAll(context.Context) ([]*biz.Feed, error) {
	return nil, nil
}
