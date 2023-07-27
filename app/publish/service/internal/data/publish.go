package data

import (
	"context"

	"Atreus/app/publish/service/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type publishRepo struct {
	data *Data
	log  *log.Helper
}

// NewPublishRepo .
func NewPublishRepo(data *Data, logger log.Logger) biz.PublishRepo {
	return &publishRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *publishRepo) Save(ctx context.Context, g *biz.Publish) (*biz.Publish, error) {
	return g, nil
}

func (r *publishRepo) Update(ctx context.Context, g *biz.Publish) (*biz.Publish, error) {
	return g, nil
}

func (r *publishRepo) FindByID(context.Context, int64) (*biz.Publish, error) {
	return nil, nil
}

func (r *publishRepo) ListByHello(context.Context, string) ([]*biz.Publish, error) {
	return nil, nil
}

func (r *publishRepo) ListAll(context.Context) ([]*biz.Publish, error) {
	return nil, nil
}
