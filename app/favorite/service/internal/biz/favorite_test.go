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
		VideoID: 1,
		UserID:  2,
	},
}
var testVideoData = map[uint32]Video{
	1: {
		Id:            1,
		Author:        &User{Id: 1},
		PlayUrl:       "https://www.baidu.com",
		CoverUrl:      "https://www.baidu.com",
		FavoriteCount: 2,
		CommentCount:  1,
		IsFavorite:    true,
		Title:         "test1",
	},
	2: {
		Id:            2,
		Author:        &User{Id: 1},
		PlayUrl:       "https://www.baidu.com",
		CoverUrl:      "https://www.baidu.com",
		FavoriteCount: 1,
		CommentCount:  1,
		IsFavorite:    false,
		Title:         "test2",
	},
	3: {
		Id:            3,
		Author:        &User{Id: 1},
		PlayUrl:       "https://www.baidu.com",
		CoverUrl:      "https://www.baidu.com",
		FavoriteCount: 1,
		CommentCount:  1,
		IsFavorite:    false,
		Title:         "test3",
	},
}
var testUserData = map[uint32]User{
	1: {
		Id:             1,
		Name:           "test1",
		FavoriteCount:  3,
		TotalFavorited: 4,
	},
	2: {
		Id:             2,
		Name:           "test2",
		FavoriteCount:  1,
		TotalFavorited: 0,
	},
}

var autoCount uint32 = 5

type MockFavoriteRepo struct{}

func (m *MockFavoriteRepo) DeleteFavoriteTx(ctx context.Context, userId, videoId uint32) error {
	for k, v := range testFavoriteData {
		if v.UserID == userId && v.VideoID == videoId {
			delete(testFavoriteData, k)
			break
		}
	}
	return nil
}

func (m *MockFavoriteRepo) CreateFavoriteTx(ctx context.Context, userId, videoId uint32) error {
	testFavoriteData[autoCount] = Favorite{
		VideoID: videoId,
		UserID:  userId,
	}
	autoCount++
	return nil
}

func (m *MockFavoriteRepo) GetFavoriteList(ctx context.Context, userId uint32) ([]Video, error) {
	var favorites []Video
	for k, v := range testFavoriteData {
		if v.UserID == userId {
			favorites = append(favorites, testVideoData[k])
		}
	}
	return favorites, nil
}
func (m *MockFavoriteRepo) IsFavorite(ctx context.Context, userId uint32, videoId []uint32) ([]bool, error) {
	isFavorite := make([]bool, len(videoId))
	for i, v := range videoId {
		for _, f := range testFavoriteData {
			if f.UserID == userId && f.VideoID == v {
				isFavorite[i] = true
				break
			}
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
		context.Background(), 1, 2, token)
	assert.Nil(t, err)
	err = usecase.FavoriteAction(
		context.Background(), 1, 1, token)
	assert.Nil(t, err)
	err = usecase.FavoriteAction(
		context.Background(), 1, 3, token)
	assert.NotEqual(t, err, nil)
}

func TestFavoriteUsecase_GetFavoriteList(t *testing.T) {
	favorites, err := usecase.GetFavoriteList(context.TODO(), 1, token)
	assert.Nil(t, err)
	var count int = 0
	for _, v := range testFavoriteData {
		if v.UserID == 1 {
			count++
		}
	}
	assert.Equal(t, len(favorites), count)
}

func TestFavoriteUsecase_IsFavorite(t *testing.T) {
	isFavorite, err := usecase.IsFavorite(context.TODO(), 1, []uint32{6})
	assert.Nil(t, err)
	assert.Equal(t, isFavorite[0], false)
	isFavorite, err = usecase.IsFavorite(context.TODO(), 1, []uint32{1})
	assert.Nil(t, err)
	assert.Equal(t, isFavorite[0], true)
}
