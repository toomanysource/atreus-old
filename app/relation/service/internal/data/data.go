package data

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"gorm.io/gorm"
)

var ProviderSet = wire.NewSet(NewData, NewRelationRepo, NewUserRepo, NewMysqlConn, NewRedisConn)

type RedisConn struct {
	client *redis.Client
	cache  *cache.Cache
}

type CacheClient struct {
	relationList *RedisConn
}

type Data struct {
	db    *gorm.DB
	cache *CacheClient
	log   *log.Helper
}

func NewData(logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{}, cleanup, nil
}
