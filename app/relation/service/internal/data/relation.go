package data

import (
	"Atreus/app/relation/service/internal/biz"
	"context"
	"google.golang.org/grpc"
	"gorm.io/gorm"

	"github.com/go-kratos/kratos/v2/log"
)

type Followers struct {
	Id         uint32 `gorm:"primary_key"`
	UserId     uint32 `gorm:"column:user_id;not null"`
	FollowerId uint32 `gorm:"column:follower_id;not null"`
	CreateAt   string `gorm:"column:created_at;default:''"`
	gorm.DeletedAt
}

func (Followers) TableName() string {
	return "followers"
}

type relationRepo struct {
	data     *Data
	userRepo *UserRepo
	log      *log.Helper
}

func NewRelationRepo(data *Data, conn *grpc.ClientConn, logger log.Logger) biz.RelationRepo {
	return &relationRepo{
		data:     data,
		userRepo: NewUserRepo(conn),
		log:      log.NewHelper(logger),
	}
}

func (r *relationRepo) GetFollowList(ctx context.Context, userId uint32) ([]*biz.User, error) {
	return r.GetFlList(ctx, userId)
}

func (r *relationRepo) GetFollowerList(ctx context.Context, userId uint32) ([]*biz.User, error) {
	return r.GetFlrList(ctx, userId)
}

func (r *relationRepo) Follow(ctx context.Context, userId uint32, toUserId uint32) error {
	return r.Follow(ctx, userId, toUserId)
}

func (r *relationRepo) UnFollow(ctx context.Context, userId uint32, toUserId uint32) error {
	return r.UnFollow(ctx, userId, toUserId)
}

func (r *relationRepo) GetFlList(ctx context.Context, userId uint32) ([]*biz.User, error) {

}

func (r *relationRepo) GetFlrList(ctx context.Context, userId uint32) ([]*biz.User, error) {

}
