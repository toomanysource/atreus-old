package biz

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/toomanysource/atreus/app/user/service/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/assert"
)

var ctx = context.Background()

type MockUserRepo struct{}

func (m *MockUserRepo) Save(ctx context.Context, user *User) (*User, error) {
	if user.Username == "foo" {
		return user, nil
	}
	return &User{}, nil
}

func (m *MockUserRepo) FindById(ctx context.Context, id uint32) (*User, error) {
	if id < 3 {
		return &User{Id: id}, nil
	}
	s := strconv.Itoa(int(id))
	return &User{Id: id, Username: s, Password: s}, nil
}

func (m *MockUserRepo) FindByIds(ctx context.Context, userId uint32, ids []uint32) ([]*User, error) {
	us := []*User{}
	for _, id := range ids {
		u, _ := m.FindById(context.Background(), id)
		if u.Username == "" {
			continue
		}
		us = append(us, u)
	}
	return us, nil
}

func (m *MockUserRepo) FindByUsername(ctx context.Context, username string) (*User, error) {
	if username == "foo" {
		return &User{}, nil
	}
	return &User{Username: username, Password: username}, nil
}

func (m *MockUserRepo) UpdateFollow(ctx context.Context, id uint32, followChange int32) error {
	return nil
}

func (m *MockUserRepo) UpdateFollower(ctx context.Context, id uint32, followerChange int32) error {
	return nil
}

func (m *MockUserRepo) UpdateFavorited(ctx context.Context, id uint32, favoritedChange int32) error {
	return nil
}

func (m *MockUserRepo) UpdateWork(ctx context.Context, id uint32, workChange int32) error {
	return nil
}

func (m *MockUserRepo) UpdateFavorite(ctx context.Context, id uint32, favoriteChange int32) error {
	return nil
}

var testConfig = &conf.JWT{
	Http: &conf.JWT_Http{
		TokenKey: "AtReUs",
	},
}
var mockRepo = &MockUserRepo{}

var usecase *UserUsecase

func TestMain(m *testing.M) {
	ctx = context.WithValue(ctx, "userId", uint32(1))
	usecase = NewUserUsecase(mockRepo, testConfig, log.DefaultLogger)
	r := m.Run()
	os.Exit(r)
}

func TestUserRegister(t *testing.T) {
	// user has been registered
	_, err := usecase.Register(ctx, "xxx", "xxx")
	assert.Error(t, err)
	// user can register
	user, err := usecase.Register(ctx, "foo", "bar")
	assert.NoError(t, err)
	assert.Equal(t, user.Username, "foo")
}

func TestUserLogin(t *testing.T) {
	// user not register
	_, err := usecase.Login(ctx, "foo", "bar")
	assert.Error(t, err)
	// incorrect password
	_, err = usecase.Login(ctx, "bar", "foo")
	assert.Error(t, err)
	// login success
	user, err := usecase.Login(ctx, "xxx", "xxx")
	assert.NoError(t, err)
	assert.Equal(t, user.Username, "xxx")
}

func TestGetInfo(t *testing.T) {
	// user not found
	_, err := usecase.GetInfo(ctx, 2)
	assert.Error(t, err)
	// user can find
	user, err := usecase.GetInfo(ctx, 4)
	assert.NoError(t, err)
	assert.Equal(t, user.Username, "4")
}

func TestGetInfos(t *testing.T) {
	// all ids can find user
	ids := []uint32{3, 4, 5, 6, 7}
	userId := uint32(1)
	users, _ := usecase.GetInfos(ctx, userId, ids)
	assert.Equal(t, len(users), len(ids))
	// some ids can not find user
	ids = []uint32{2, 3, 4, 5, 6}
	users, _ = usecase.GetInfos(ctx, userId, ids)
	assert.Equal(t, len(users), len(ids)-1)
}

func TestUpdate(t *testing.T) {
	err := usecase.UpdateFollow(ctx, 1, 1)
	assert.NoError(t, err)
	err = usecase.UpdateFollower(ctx, 2, 2)
	assert.NoError(t, err)
	err = usecase.UpdateFavorite(ctx, 3, 3)
	assert.NoError(t, err)
	err = usecase.UpdateFavorited(ctx, 4, 4)
	assert.NoError(t, err)
	err = usecase.UpdateWork(ctx, 5, 5)
	assert.NoError(t, err)
}
