package data

import (
	"context"

	"Atreus/app/favorite/service/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type favoriteRepo struct {
	data *Data
	log  *log.Helper
}

// NewFavoriteRepo .
func NewFavoriteRepo(data *Data, logger log.Logger) biz.FavoriteRepo {
	return &favoriteRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *favoriteRepo) Save(ctx context.Context, g *biz.Favorite) (*biz.Favorite, error) {
	return g, nil
}

func (r *favoriteRepo) Update(ctx context.Context, g *biz.Favorite) (*biz.Favorite, error) {
	return g, nil
}

func (r *favoriteRepo) FindByID(context.Context, int64) (*biz.Favorite, error) {
	return nil, nil
}

func (r *favoriteRepo) ListByHello(context.Context, string) ([]*biz.Favorite, error) {
	return nil, nil
}

func (r *favoriteRepo) ListAll(context.Context) ([]*biz.Favorite, error) {
	return nil, nil
}
