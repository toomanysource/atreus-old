package biz

import (
	"context"
	"errors"

	"github.com/go-kratos/kratos/v2/log"
)

// var (
// 	// ErrUserNotFound is user not found.
// 	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
// )

// User is a user model.
type User struct {
	ID              int64
	Username        string
	Password        string
	Name            string
	FollowCount     int64
	FollowerCount   int64
	IsFollow        bool
	Avatar          string
	BackgroundImage string
	Signature       string
	TotalFavorite   int64
	WorkCount       int64
	FavoriteCount   int64
}

// UserRepo is a user repo.
type UserRepo interface {
	Save(context.Context, *User) (*User, error)
	FindByID(context.Context, int64) (*User, error)
	FindByUsername(context.Context, string) (*User, error)
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
	if user.Username == "" {
		return nil, errors.New("the username has been registered")
	}

	regUser := &User{
		Username: username,
		Password: password,
	}
	user, err = uc.repo.Save(ctx, regUser)
	if err != nil {
		return nil, errors.New("register service not work")
	}
	return user, nil
}

// Login .
func (uc *UserUsecase) Login(ctx context.Context, username, password string) (*User, error) {
	user, _ := uc.repo.FindByUsername(ctx, username)
	if user.Username == "" {
		return nil, errors.New("can not find registered user")
	}

	if user.Password != password {
		return nil, errors.New("incorrect password")
	}

	return user, nil
}

// GetInfo .
func (uc *UserUsecase) GetInfo(ctx context.Context, userId int64) (*User, error) {
	user, _ := uc.repo.FindByID(ctx, userId)
	if user.Username == "" {
		return nil, errors.New("can not find the user")
	}

	return user, nil
}
