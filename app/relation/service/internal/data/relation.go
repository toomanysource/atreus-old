package data

import (
	"context"

	"Atreus/app/relation/service/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type relationRepo struct {
	data *Data
	log  *log.Helper
}

// NewRelationRepo .
func NewRelationRepo(data *Data, logger log.Logger) biz.RelationRepo {
	return &relationRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *relationRepo) Save(ctx context.Context, g *biz.Relation) (*biz.Relation, error) {
	return g, nil
}

func (r *relationRepo) Update(ctx context.Context, g *biz.Relation) (*biz.Relation, error) {
	return g, nil
}

func (r *relationRepo) FindByID(context.Context, int64) (*biz.Relation, error) {
	return nil, nil
}

func (r *relationRepo) ListByHello(context.Context, string) ([]*biz.Relation, error) {
	return nil, nil
}

func (r *relationRepo) ListAll(context.Context) ([]*biz.Relation, error) {
	return nil, nil
}
