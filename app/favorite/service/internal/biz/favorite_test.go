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
		FavoriteCount: 1,
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
		TotalFavorited: 3,
	},
	2: {
		Id:             2,
		Name:           "test2",
		FavoriteCount:  0,
		TotalFavorited: 1,
	},
}
var autoCount uint32 = 5

type MockFavoriteRepo struct{}

func (m *MockFavoriteRepo) DeleteFavoriteTx(ctx context.Context, userId, videoId, authorId uint32) error {
	err := m.DeleteFavorite(ctx, videoId, userId)
	if err != nil {
		return err
	}
	err = userRepo.UpdateFavorite(ctx, userId, -1)
	if err != nil {
		return err
	}
	err = userRepo.UpdateFavorited(ctx, authorId, -1)
	if err != nil {
		return err
	}
	return nil
}

func (m *MockFavoriteRepo) CreateFavoriteTx(ctx context.Context, userId, videoId, authorId uint32) error {
	err := m.CreateFavorite(ctx, videoId, userId)
	if err != nil {
		return err
	}
	err = userRepo.UpdateFavorite(ctx, userId, 1)
	if err != nil {
		return err
	}
	err = userRepo.UpdateFavorited(ctx, authorId, 1)
	if err != nil {
		return err
	}
	return nil
}

func (m *MockFavoriteRepo) CreateFavorite(ctx context.Context, videoId, userId uint32) error {
	favorite := Favorite{
		VideoID: videoId,
		UserID:  userId,
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
func (m *MockFavoriteRepo) IsFavorite(ctx context.Context, videoId uint32, userId uint32) (bool, error) {
	//return int64(len(testFavoriteData)), nil
	isFavorite := false
	for _, v := range testFavoriteData {
		if videoId == v.VideoID && v.UserID == userId {
			isFavorite = true
			break
		}
	}
	return isFavorite, nil
}

type MockUserRepo struct{}

func (m *MockUserRepo) UpdateFavorite(ctx context.Context, userId uint32, change int32) error {
	//testUserData[userId].FavoriteCount = testUserData[userId].FavoriteCount + uint32(change)
	user := testUserData[userId]
	user.FavoriteCount = user.FavoriteCount + uint32(change)
	testUserData[userId] = user
	return nil
}
func (m *MockUserRepo) UpdateFavorited(ctx context.Context, userId uint32, change int32) error {
	return nil
}

type MockPublishRepo struct{}
type MockTransaction struct{}

func (m *MockPublishRepo) GetVideoListByVideoIds(ctx context.Context, videoIds []uint32) ([]Video, error) {
	var videos []Video
	for _, v := range videoIds {
		videos = append(videos, testVideoData[v])
	}
	return videos, nil
}

func (m *MockTransaction) ExecTx(context.Context, func(ctx context.Context) error) error {
	return nil
}

var mockRepo = &MockFavoriteRepo{}
var userRepo = &MockUserRepo{}
var publishRepo = &MockPublishRepo{}
var mockTransaction = &MockTransaction{}

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
	usecase = NewFavoriteUsecase(testConfig, mockRepo, userRepo, publishRepo, mockTransaction, log.DefaultLogger)
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
	assert.Equal(t, isFavorite, false)
	isFavorite, err = usecase.IsFavorite(context.TODO(), 1, []uint32{1})
	assert.Nil(t, err)
	assert.Equal(t, isFavorite, true)
}
