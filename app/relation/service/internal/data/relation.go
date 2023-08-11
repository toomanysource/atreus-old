package data

import (
	"Atreus/app/comment/service/pkg/gormX"
	"Atreus/app/relation/service/internal/biz"
	"context"
	"google.golang.org/grpc"
	"gorm.io/gorm"

	"github.com/go-kratos/kratos/v2/log"
)

type UserRepo interface {
	GetUserInfos(ctx context.Context, userIds []uint32) ([]*biz.User, error)
	UpdateFollow(ctx context.Context, userId uint32, followChange int32) error
	UpdateFollower(ctx context.Context, userId uint32, followerChange int32) error
}

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
	userRepo UserRepo
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

func (r *relationRepo) IsFollow(ctx context.Context, userId uint32, toUserId uint32) (bool, error) {
	relation, err := r.SearchRelation(ctx, userId, toUserId)
	if err != nil {
		return false, err
	}
	if relation == nil {
		return false, nil
	}
	return true, nil
}

// GetFlList 获取关注列表
func (r *relationRepo) GetFlList(ctx context.Context, userId uint32) (users []*biz.User, err error) {
	var follows []*Followers
	err = r.data.db.Action(ctx, func(tran gormX.Transactor) error {
		if err := tran.Where("follower_id = ?", userId).Find(&follows).Error; err != nil {
			return err
		}
		var userIDs []uint32
		for _, follow := range follows {
			userIDs = append(userIDs, follow.UserId)
		}
		users, err = r.userRepo.GetUserInfos(ctx, userIDs)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetFlrList 获取粉丝列表
func (r *relationRepo) GetFlrList(ctx context.Context, userId uint32) (users []*biz.User, err error) {
	var followers []*Followers
	err = r.data.db.Action(ctx, func(tran gormX.Transactor) error {
		if err = tran.Where("user_id = ?", userId).Find(&followers).Error; err != nil {
			return err
		}
		var userIDs []uint32
		for _, follower := range followers {
			userIDs = append(userIDs, follower.FollowerId)
		}
		users, err = r.userRepo.GetUserInfos(ctx, userIDs)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return users, nil
}

// AddFollow 添加关注
func (r *relationRepo) AddFollow(ctx context.Context, userId uint32, toUserId uint32) error {
	err := r.data.db.Action(ctx, func(tran gormX.Transactor) error {
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
		err = tran.Create(follow).Error
		if err != nil {
			return err
		}
		err = r.userRepo.UpdateFollow(ctx, userId, 1)
		if err != nil {
			return err
		}
		err = r.userRepo.UpdateFollower(ctx, toUserId, 1)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

// DelFollow 取消关注
func (r *relationRepo) DelFollow(ctx context.Context, userId uint32, toUserId uint32) error {
	err := r.data.db.Action(ctx, func(tran gormX.Transactor) error {
		relation, err := r.SearchRelation(ctx, userId, toUserId)
		if err != nil {
			return err
		}
		err = tran.Delete(relation).Error
		if err != nil {
			return err
		}
		err = r.userRepo.UpdateFollow(ctx, userId, -1)
		if err != nil {
			return err
		}
		err = r.userRepo.UpdateFollower(ctx, toUserId, -1)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

// SearchRelation 查询关注关系
func (r *relationRepo) SearchRelation(ctx context.Context, userId uint32, toUserId uint32) (*Followers, error) {
	var relation *Followers
	err := r.data.db.Action(ctx, func(tran gormX.Transactor) error {
		if err := tran.WithContext(ctx).Where(
			"user_id = ? and follower_id = ?", userId, toUserId).Find(relation).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return relation, nil
}
