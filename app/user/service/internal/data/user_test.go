// Before testing this file, you should make sure the database has already running.
// If you didn't use the file atreus/deploy/dockercompose/docker-compose.yaml to start mysql service,
// you should check the source of database and modify it if it's necessary.
package data

import (
	"math/rand"
	"os"
	"testing"

	"github.com/toomanysource/atreus/app/user/service/internal/biz"
	"github.com/toomanysource/atreus/app/user/service/internal/conf"
	"github.com/toomanysource/atreus/app/user/service/internal/server"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

var (
	ctx           = context.Background()
	testUsersData = []*biz.User{
		{
			Id:       1,
			Username: "uuuuuusername114",
			Password: "ppppppassword114",
		},
		{
			Id:       2,
			Username: "foo114",
			Password: "bar114",
		},
		{
			Id:       3,
			Username: "OuO114",
			Password: "NuN114",
		},
		{
			Id:       4,
			Username: "okgogogo114",
			Password: "okgogogo114",
		},
		{
			Id:       5,
			Username: "someone114",
			Password: "another114",
		},
		{
			Id:       6,
			Username: "i_dont_care_what_username_is",
			Password: "i_dont_care_what_password_is",
		},
	}
)

var testConfig = &conf.Data{
	Database: &conf.Data_Database{
		Driver: "mysql",
		// if you don't use default config, the source must be modified
		Source: "root:atreus114@tcp(127.0.0.1:33069)/atreus?charset=utf8mb4&parseTime=True&loc=Local",
	},
}

var relationConfig = &conf.Client{
	Relation: &conf.Client_Relation{
		To: "0.0.0.0:9004",
	},
}

var repo *userRepo

func TestMain(m *testing.M) {
	ctx = context.WithValue(ctx, "user_id", uint32(1))
	db := NewGormDb(testConfig)
	logger := log.DefaultLogger
	relationConn := server.NewRelationClient(relationConfig, logger)
	data, f, err := NewData(db, logger)
	if err != nil {
		panic(err)
	}
	repo = (NewUserRepo(data, relationConn, logger)).(*userRepo)
	r := m.Run()
	f()
	os.Exit(r)
}

func TestSaveUser(t *testing.T) {
	for _, user := range testUsersData {
		_, err := repo.Save(ctx, user)
		assert.NoError(t, err)
	}
}

func TestFindUserById(t *testing.T) {
	for _, user := range testUsersData {
		_, err := repo.FindById(ctx, user.Id)
		assert.NoError(t, err)
	}
}

func TestFindUserByIds(t *testing.T) {
	userIds := make([]uint32, len(testUsersData))
	for i, user := range testUsersData {
		userIds[i] = user.Id
	}
	var userId uint32 = 1
	users, err := repo.FindByIds(ctx, userId, userIds)
	assert.NoError(t, err)
	assert.Equal(t, len(userIds), len(users))
}

func TestFindByUsername(t *testing.T) {
	for _, user := range testUsersData {
		u, err := repo.FindByUsername(ctx, user.Username)
		assert.NoError(t, err)
		assert.Equal(t, user.Username, u.Username)
	}
}

func TestUpdateUser(t *testing.T) {
	id := testUsersData[0].Id
	var rangNum int32 = 20
	t.Run("updateUserFavorite", func(t *testing.T) {
		user, err := repo.FindById(ctx, id)
		assert.Nil(t, err)
		old := user.FavoriteCount

		change := rand.Int31n(rangNum)
		err = repo.UpdateFavorite(ctx, id, change)
		assert.Nil(t, err)
		user, _ = repo.FindById(ctx, id)
		assert.Equal(t, old+uint32(change), user.FavoriteCount)
	})
	t.Run("updateUserFavorited", func(t *testing.T) {
		user, err := repo.FindById(ctx, id)
		assert.Nil(t, err)
		old := user.TotalFavorited

		change := rand.Int31n(rangNum)
		err = repo.UpdateFavorited(ctx, id, change)
		assert.Nil(t, err)
		user, _ = repo.FindById(ctx, id)
		assert.Equal(t, old+uint32(change), user.TotalFavorited)
	})
	t.Run("updateUserFollow", func(t *testing.T) {
		user, err := repo.FindById(ctx, id)
		assert.Nil(t, err)
		old := user.FollowCount

		change := rand.Int31n(rangNum)
		err = repo.UpdateFollow(ctx, id, change)
		assert.Nil(t, err)
		user, _ = repo.FindById(ctx, id)
		assert.Equal(t, old+uint32(change), user.FollowCount)
	})
	t.Run("updateUserFollower", func(t *testing.T) {
		user, err := repo.FindById(ctx, id)
		assert.Nil(t, err)
		old := user.FollowerCount

		change := rand.Int31n(rangNum)
		err = repo.UpdateFollower(ctx, id, change)
		assert.Nil(t, err)
		user, _ = repo.FindById(ctx, id)
		assert.Equal(t, old+uint32(change), user.FollowerCount)
	})
	t.Run("updateUserWork", func(t *testing.T) {
		user, err := repo.FindById(ctx, id)
		assert.Nil(t, err)
		old := user.WorkCount

		change := rand.Int31n(rangNum)
		err = repo.UpdateWork(ctx, id, change)
		assert.Nil(t, err)
		user, _ = repo.FindById(ctx, id)
		assert.Equal(t, old+uint32(change), user.WorkCount)
	})
}
