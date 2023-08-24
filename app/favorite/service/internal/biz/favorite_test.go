package biz

import (
	"context"
	"os"
	"testing"

	"Atreus/app/favorite/service/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

var testVideoData = []Video{
	{
		Id:            1,
		Author:        &User{Id: 1},
		PlayUrl:       "https://www.baidu.com",
		CoverUrl:      "https://www.baidu.com",
		FavoriteCount: 2,
		CommentCount:  1,
		IsFavorite:    true,
		Title:         "test1",
	},
	{
		Id:            2,
		Author:        &User{Id: 1},
		PlayUrl:       "https://www.baidu.com",
		CoverUrl:      "https://www.baidu.com",
		FavoriteCount: 1,
		CommentCount:  1,
		IsFavorite:    false,
		Title:         "test2",
	},
	{
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

type MockFavoriteRepo struct{}

func (m *MockFavoriteRepo) DeleteFavorite(ctx context.Context, userId, videoId uint32) error {
	return nil
}

func (m *MockFavoriteRepo) CreateFavorite(ctx context.Context, userId, videoId uint32) error {
	return nil
}

func (m *MockFavoriteRepo) GetFavoriteList(ctx context.Context, userId uint32) ([]Video, error) {
	var favoriteList []Video
	for _, v := range testVideoData {
		if v.IsFavorite == true {
			favoriteList = append(favoriteList, v)
		}
	}
	return favoriteList, nil
}

func (m *MockFavoriteRepo) IsFavorite(ctx context.Context, userId uint32, videoId []uint32) ([]bool, error) {
	isFavorite := make([]bool, len(videoId))
	for i := range videoId {
		isFavorite[i] = false
	}
	return isFavorite, nil
}

var (
	mockRepo   = &MockFavoriteRepo{}
	usecase    *FavoriteUsecase
	testConfig = &conf.JWT{
		Http: &conf.JWT_HTTP{
			TokenKey: "TEST",
		},
	}
)

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
	for _, v := range favorites {
		assert.Equal(t, v.IsFavorite, true)
	}
}

func TestFavoriteUsecase_IsFavorite(t *testing.T) {
	isFavorite, err := usecase.IsFavorite(context.TODO(), 1, []uint32{6})
	assert.Nil(t, err)
	assert.Equal(t, isFavorite[0], false)
	isFavorite, err = usecase.IsFavorite(context.TODO(), 1, []uint32{1})
	assert.Nil(t, err)
	assert.Equal(t, isFavorite[0], true)
}
