package biz

import (
	"Atreus/app/favorite/service/internal/conf"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var testFavoriteData = map[uint32]Favorite{
	1: {
		VideoID: 1,
		UserID:  1,
	},
	2: {
		VideoID: 2,
		UserID:  1,
	},
	3: {
		VideoID: 3,
		UserID:  1,
	},
	4: {
		VideoID: 4,
		UserID:  1,
	},
	5: {
		VideoID: 5,
		UserID:  1,
	},
	6: {
		VideoID: 6,
		UserID:  1,
	},
}
var testVideoData = map[uint32]Video{
	1: {
		Id:            1,
		Author:        &User{},
		PlayUrl:       "https://www.baidu.com",
		CoverUrl:      "https://www.baidu.com",
		FavoriteCount: 1,
		CommentCount:  1,
		IsFavorite:    true,
		Title:         "test1",
	},
	2: {
		Id:            2,
		Author:        &User{},
		PlayUrl:       "https://www.baidu.com",
		CoverUrl:      "https://www.baidu.com",
		FavoriteCount: 1,
		CommentCount:  1,
		IsFavorite:    false,
		Title:         "test2",
	},
	3: {
		Id:            3,
		Author:        &User{},
		PlayUrl:       "https://www.baidu.com",
		CoverUrl:      "https://www.baidu.com",
		FavoriteCount: 1,
		CommentCount:  1,
		IsFavorite:    false,
		Title:         "test3",
	},
}

var autoCount uint32 = 7

type MockFavoriteRepo struct{}

func (m *MockFavoriteRepo) CreateFavorite(ctx context.Context, videoId, userId uint32) error {
	favorite := Favorite{
		videoId, userId,
	}
	testFavoriteData[autoCount] = favorite
	autoCount++
	return nil
}

func (m *MockFavoriteRepo) DeleteFavorite(ctx context.Context, videoId, userId uint32) error {
	for k, v := range testFavoriteData {
		if v.VideoID == videoId && v.UserID == userId {
			delete(testFavoriteData, k)
		}
	}
	return nil
}

func (m *MockFavoriteRepo) GetFavoriteList(ctx context.Context, userId uint32) ([]Video, error) {
	var favorites []Video
	for k, v := range testFavoriteData {
		//favorites = append(favorites, favorite)
		if v.UserID == userId {
			favorites = append(favorites, testVideoData[k])
		}
	}
	return favorites, nil
}
func (m *MockFavoriteRepo) IsFavorite(ctx context.Context, videoId, userId uint32) (bool, error) {
	//return int64(len(testFavoriteData)), nil
	isFavorite := false
	for _, v := range testFavoriteData {
		if v.VideoID == videoId && v.UserID == userId {
			isFavorite = true
		}
	}
	return isFavorite, nil
}

var mockRepo = &MockFavoriteRepo{}

var usecase *FavoriteUsecase

var testConfig = &conf.JWT{
	Http: &conf.JWT_HTTP{
		TokenKey: "TEST",
	},
}

var token = func() string {
	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": 1,
	}).SignedString([]byte("TEST"))
	return token
}()

func TestMain(m *testing.M) {
	usecase = NewFavoriteUsecase(testConfig, mockRepo, log.DefaultLogger)
	r := m.Run()
	os.Exit(r)
}

func TestFavoriteUsecase_FavoriteAction(t *testing.T) {
	err := usecase.FavoriteAction(
		context.TODO(), 1, 2, token)
	assert.Nil(t, err)
	err = usecase.FavoriteAction(
		context.TODO(), 1, 1, token)
	assert.Nil(t, err)
}

func TestFavoriteUsecase_GetFavoriteList(t *testing.T) {
	favorites, err := usecase.GetFavoriteList(context.TODO(), 1, token)
	assert.Nil(t, err)
	assert.Equal(t, len(favorites), len(testFavoriteData))
}

func TestFavoriteUsecase_IsFavorite(t *testing.T) {
	isFavorite, err := usecase.IsFavorite(context.TODO(), 1, 1)
	assert.Nil(t, err)
	assert.Equal(t, isFavorite, false)
	isFavorite, err = usecase.IsFavorite(context.TODO(), 1, 1)
	assert.Nil(t, err)
	assert.Equal(t, isFavorite, true)
}
