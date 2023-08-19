package data

import (
	"Atreus/app/comment/service/internal/biz"
	"Atreus/app/comment/service/internal/conf"
	"Atreus/app/comment/service/internal/server"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/durationpb"
	"os"
	"testing"
	"time"
)

var testUsersData = []*biz.User{
	{
		Id:   1,
		Name: "hahah",
	},
	{
		Id:   2,
		Name: "sefafa",
	},
	{
		Id:   3,
		Name: "brbs",
	},
	{
		Id:   4,
		Name: "awfawfaw4rt",
	},
	{
		Id:   5,
		Name: "bgssev",
	},
	{
		Id:   6,
		Name: "dawfawf",
	},
}
var testCommentsData = []*biz.Comment{
	{
		Id: 1,
		User: &biz.User{
			Id:   1,
			Name: "hahah",
		},
		Content:    "bushuwu1",
		CreateDate: "08-01",
	},
	{
		Id: 2,
		User: &biz.User{
			Id:   1,
			Name: "hahah",
		},
		Content:    "dadawd",
		CreateDate: "08-02",
	},
	{
		Id: 3,
		User: &biz.User{
			Id:   2,
			Name: "sefafa",
		},
		Content:    "bdzxvzad",
		CreateDate: "08-03",
	},
	{
		Id: 4,
		User: &biz.User{
			Id:   1,
			Name: "hahah",
		},
		Content:    "bvrbr",
		CreateDate: "08-03",
	},
	{
		Id: 5,
		User: &biz.User{
			Id:   3,
			Name: "brbs",
		},
		Content:    "bdadawfvrd",
		CreateDate: "08-04",
	},
	{
		Id: 6,
		User: &biz.User{
			Id:   5,
			Name: "bgssev",
		},
		Content:    "bdafagaagaga",
		CreateDate: "08-05",
	},
}

var testCommentsCache = map[string]map[string]string{
	"1": {
		"1": "{\"Id\": 1,\"User\":{\"Id\":1,\"Name\":\"hahah\"},\"Content\":\"bushuwu1\",\"CreateDate\":\"08-01\"}",
		"2": "{\"Id\": 2,\"User\":{\"Id\":1,\"Name\":\"hahah\"},\"Content\":\"dadawd\",\"CreateDate\":\"08-02\"}",
		"3": "{\"Id\": 3,\"User\":{\"Id\":2,\"Name\":\"sefafa\"},\"Content\":\"bdzxvzad\",\"CreateDate\":\"08-03\"}",
		"4": "{\"Id\": 4,\"User\":{\"Id\":1,\"Name\":\"hahah\"},\"Content\":\"bvrbr\",\"CreateDate\":\"08-03\"}",
		"5": "{\"Id\": 5,\"User\":{\"Id\":3,\"Name\":\"brbs\"},\"Content\":\"bdadawfvrd\",\"CreateDate\":\"08-04\"}",
	},
	"2": {
		"6": "{\"Id\": 6,\"User\":{\"Id\":5,\"Name\":\"bgssev\"},\"Content\":\"bdafagaagaga\",\"CreateDate\":\"08-05\"}",
	},
}

var testConfig = &conf.Data{
	Mysql: &conf.Data_Mysql{
		Driver: "mysql",
		// if you don't use default config, the source must be modified
		Dsn: "root:toomanysource@tcp(127.0.0.1:3306)/atreus?charset=utf8mb4&parseTime=True&loc=Local",
	},
	Redis: &conf.Data_Redis{
		CommentDb:    1,
		Addr:         "127.0.0.1:6379",
		Password:     "atreus",
		ReadTimeout:  &durationpb.Duration{Seconds: 1},
		WriteTimeout: &durationpb.Duration{Seconds: 1},
	},
}
var testClientConfig = &conf.Client{
	User: &conf.Client_User{
		To: "0.0.0.0:9001",
	},
}

var cRepo *commentRepo

func TestMain(m *testing.M) {
	logger := log.DefaultLogger
	db := NewMysqlConn(testConfig, logger)
	cache := NewRedisConn(testConfig, logger)
	userConn := server.NewUserClient(testClientConfig, logger)
	publishConn := server.NewPublishClient(testClientConfig, logger)
	data, f, err := NewData(db, cache, logger)
	if err != nil {
		panic(err)
	}
	cRepo = (NewCommentRepo(data, userConn, publishConn, logger)).(*commentRepo)
	r := m.Run()
	time.Sleep(time.Second * 2)
	f()
	os.Exit(r)
}

func TestCommentRepo_SearchCommentList(t *testing.T) {
	comments, err := cRepo.SearchCommentList(context.TODO(), 1, 1)
	assert.Nil(t, err)
	assert.Equal(t, len(comments), len(testCommentsData)-1)
}

func TestCommentRepo_InsertComment(t *testing.T) {
	_, err := cRepo.InsertComment(context.TODO(), 2, "wuhu", 2)
	assert.Nil(t, err)
}

func TestCommentRepo_DelComment(t *testing.T) {
	_, err := cRepo.DelComment(context.TODO(), 2, 19, 2)
	assert.Nil(t, err)
}

func TestCommentRepo_CacheCreateCommentTransaction(t *testing.T) {
	err := cRepo.CacheCreateCommentTransaction(context.TODO(), testCommentsData[:5], 1)
	assert.Nil(t, err)
}

func TestCommentRepo_DeleteComment(t *testing.T) {
	_, err := cRepo.DeleteComment(context.TODO(), 2, 6, 5)
	assert.Nil(t, err)
}

func TestCommentRepo_CreateComment(t *testing.T) {
	_, err := cRepo.CreateComment(context.TODO(), 2, "hahaha", 3)
	assert.Nil(t, err)
}

func TestCommentRepo_GetCommentList(t *testing.T) {
	comments, err := cRepo.GetCommentList(context.TODO(), 1, 2)
	assert.Nil(t, err)
	assert.Equal(t, len(comments), 5)
}
