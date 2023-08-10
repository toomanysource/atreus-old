package biz

import (
	"Atreus/app/user/service/internal/pkg"
	"Atreus/pkg/common"
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

var ErrInternal = errors.New("internal error")

// User is a user model.
type User struct {
	Id              uint32    `gorm:"primary_key"`
	Username        string    `gorm:"column:username;not null"`
	Password        string    `gorm:"column:password;not null"`
	Name            string    `gorm:"column:name;not null"`
	FollowCount     uint32    `gorm:"column:follow_count;not null;default:0"`
	FollowerCount   uint32    `gorm:"column:follower_count;not null;default:0"`
	Avatar          string    `gorm:"column:avatar_url;not null;default:''"`
	BackgroundImage string    `gorm:"column:background_image_url;not null;default:''"`
	Signature       string    `gorm:"column:signature;not null;default:''"`
	TotalFavorited  uint32    `gorm:"column:total_favorited;not null;default:0"`
	WorkCount       uint32    `gorm:"column:work_count;not null;default:0"`
	FavoriteCount   uint32    `grom:"column:favorite_count;not null;default:0"`
	IsFollow        bool      `gorm:"-"`
	Salt            string    `gorm:"column:salt"`
	Created_at      time.Time `gorm:"column:created_at"`
	Updated_at      time.Time `gorm:"column:updated_at"`
	Token           string    `gorm:"-"`
}

// UserInfo is the information that user can modify
type UserInfo struct {
	Name            string
	Avatar          string
	BackgroundImage string
	Signature       string
}

// UserRepo is a user repo.
type UserRepo interface {
	Save(context.Context, *User) (*User, error)
	FindById(context.Context, uint32) (*User, error)
	FindByIds(context.Context, []uint32) ([]*User, error)
	FindByUsername(context.Context, string) (*User, error)
	UpdateInfo(context.Context, *UserInfo) error
	UpdateFollow(context.Context, uint32, int32) error
	UpdateFollower(context.Context, uint32, int32) error
	UpdateFavorited(context.Context, uint32, int32) error
	UpdateWork(context.Context, uint32, int32) error
	UpdateFavorite(context.Context, uint32, int32) error
}

// UserUsecase is a user usecase.
type UserUsecase struct {
	repo UserRepo
	log  *log.Helper
}

// NewUserUsecase new a user usecase.
func NewUserUsecase(repo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{repo: repo, log: log.NewHelper(logger)}
}

// Register .
func (uc *UserUsecase) Register(ctx context.Context, username, password string) (*User, error) {
	user, err := uc.repo.FindByUsername(ctx, username)
	if err != nil {
		return nil, ErrInternal
	}
	if user.Username != "" {
		return nil, errors.New("the username has been registered")
	}

	salt := pkg.RandomString(10)
	password = pkg.GenSaltPassword(salt, password)

	// save user
	regUser := &User{
		Username: username,
		Password: password,
		// Name is same as username
		Name:       username,
		Salt:       salt,
		Created_at: time.Now(),
		Updated_at: time.Now(),
	}
	user, err = uc.repo.Save(ctx, regUser)
	if err != nil {
		return nil, ErrInternal
	}

	// 生成 token
	token, err := pkg.ProduceToken(user.Id)
	if err != nil {
		return nil, ErrInternal
	}
	user.Token = token
	return user, nil
}

// Login .
func (uc *UserUsecase) Login(ctx context.Context, username, password string) (*User, error) {
	user, err := uc.repo.FindByUsername(ctx, username)
	if err != nil {
		return nil, ErrInternal
	}
	if user.Username == "" {
		return nil, errors.New("can not find registered user")
	}
	password = pkg.GenSaltPassword(user.Salt, password)
	if user.Password != password {
		return nil, errors.New("incorrect password")
	}

	// 生成 token
	token, err := pkg.ProduceToken(user.Id)
	if err != nil {
		return nil, ErrInternal
	}
	user.Token = token
	return user, nil
}

// GetInfo .
func (uc *UserUsecase) GetInfo(ctx context.Context, userId uint32, tokenString string) (*User, error) {
	token, err := common.ParseToken("AtReUs", tokenString)
	if err != nil {
		return nil, err
	}
	_, err = common.GetTokenData(token)
	user, err := uc.repo.FindById(ctx, userId)
	if err != nil {
		return nil, ErrInternal
	}
	if user.Username == "" {
		return nil, errors.New("can not find the user")
	}

	return user, nil
}

// UpdateInfo not implement yet
func (uc *UserUsecase) UpdateInfo(ctx context.Context, info *UserInfo) error {
	err := uc.repo.UpdateInfo(ctx, info)
	if err != nil {
		return ErrInternal
	}
	return nil
}

// GetInfos .
func (uc *UserUsecase) GetInfos(ctx context.Context, userIds []uint32) ([]*User, error) {
	users, err := uc.repo.FindByIds(ctx, userIds)
	if err != nil {
		return nil, ErrInternal
	}
	if len(users) == 0 {
		return []*User{}, nil
	}

	return users, nil
}

// UpdateFollow .
func (uc *UserUsecase) UpdateFollow(ctx context.Context, userId uint32, followChange int32) error {
	err := uc.repo.UpdateFollow(ctx, userId, followChange)
	if err != nil {
		return ErrInternal
	}

	return nil
}

// UpdateFollower .
func (uc *UserUsecase) UpdateFollower(ctx context.Context, userId uint32, followerChange int32) error {
	err := uc.repo.UpdateFollower(ctx, userId, followerChange)
	if err != nil {
		return ErrInternal
	}

	return nil
}

// UpdateFavorited .
func (uc *UserUsecase) UpdateFavorited(ctx context.Context, userId uint32, favoritedChange int32) error {
	err := uc.repo.UpdateFavorited(ctx, userId, favoritedChange)
	if err != nil {
		return ErrInternal
	}

	return nil
}

// UpdateWork .
func (uc *UserUsecase) UpdateWork(ctx context.Context, userId uint32, workChange int32) error {
	err := uc.repo.UpdateWork(ctx, userId, workChange)
	if err != nil {
		return ErrInternal
	}

	return nil
}

// UpdateFavorite .
func (uc *UserUsecase) UpdateFavorite(ctx context.Context, userId uint32, favoriteChange int32) error {
	err := uc.repo.UpdateFavorite(ctx, userId, favoriteChange)
	if err != nil {
		return ErrInternal
	}

	return nil
}
