package data

import (
	"Atreus/app/relation/service/internal/biz"
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type relationRepo struct {
	data *Data
	log  *log.Helper
}

func NewRelationRepo(data *Data, logger log.Logger) biz.RelationRepo {
	return &relationRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *relationRepo) GetFollowList(ctx context.Context, userId int64) ([]*biz.User, error) {

	return []*biz.User{}, nil
}

func (r *relationRepo) GetFollowerList(ctx context.Context, userId int64) ([]*biz.User, error) {
	return []*biz.User{}, nil
}

func (r *relationRepo) Follow(ctx context.Context, userId int64, toUserId int64) error {
	return nil
}

func (r *relationRepo) UnFollow(ctx context.Context, userId int64, toUserId int64) error {
	return nil
}
