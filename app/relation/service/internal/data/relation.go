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
	return r.AddFollow(ctx, userId, toUserId)
}

func (r *relationRepo) UnFollow(ctx context.Context, userId uint32, toUserId uint32) error {
	return r.DelFollow(ctx, userId, toUserId)
}

// GetFlList 获取关注列表
func (r *relationRepo) GetFlList(ctx context.Context, userId uint32) ([]*biz.User, error) {
	var follows []*Followers
	if err := r.data.db.Where("follower_id = ?", userId).Find(&follows).Error; err != nil {
		return nil, err
	}
	var userIDs []uint32
	for _, follow := range follows {
		userIDs = append(userIDs, follow.UserId)
	}
	users, err := r.userRepo.GetUserInfos(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetFlrList 获取粉丝列表
func (r *relationRepo) GetFlrList(ctx context.Context, userId uint32) ([]*biz.User, error) {
	var followers []*Followers
	if err := r.data.db.Where("user_id = ?", userId).Find(&followers).Error; err != nil {
		return nil, err
	}
	var userIDs []uint32
	for _, follower := range followers {
		userIDs = append(userIDs, follower.FollowerId)
	}
	users, err := r.userRepo.GetUserInfos(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// AddFollow 添加关注
func (r *relationRepo) AddFollow(ctx context.Context, userId uint32, toUserId uint32) error {
	relation, err := r.SearchRelation(ctx, userId, toUserId)
	if err != nil {
		return err
	}
	if relation != nil {
		return nil
	}
	follow := &Followers{
		UserId:     userId,
		FollowerId: toUserId,
	}
	return r.data.db.Create(follow).Error
}

// DelFollow 取消关注
func (r *relationRepo) DelFollow(ctx context.Context, userId uint32, toUserId uint32) error {
	relation, err := r.SearchRelation(ctx, userId, toUserId)
	if err != nil {
		return err
	}
	return r.data.db.Delete(relation).Error
}

// GetFollowCount 获取关注数
func (r *relationRepo) GetFollowCount(ctx context.Context, userId uint32) (int64, error) {
	var count int64
	if err := r.data.db.Where("follower_id = ?", userId).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetFollowerCount 获取粉丝数
func (r *relationRepo) GetFollowerCount(ctx context.Context, userId uint32) (int64, error) {
	var count int64
	if err := r.data.db.Where("user_id = ?", userId).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// SearchRelation 查询关注关系
func (r *relationRepo) SearchRelation(ctx context.Context, userId uint32, toUserId uint32) (*Followers, error) {
	var relation *Followers
	if err := r.data.db.Where(
		"user_id = ? and follower_id = ?", userId, toUserId).Find(relation).Error; err != nil {
		return nil, err
	}
	return relation, nil
}
