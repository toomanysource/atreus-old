package data

import (
	"Atreus/app/relation/service/internal/conf"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sync"
)

var ProviderSet = wire.NewSet(NewData, NewRelationRepo, NewUserRepo, NewMysqlConn, NewRedisConn)

type Data struct {
	db    *gorm.DB
	cache *redis.Client
	log   *log.Helper
}

func NewData(db *gorm.DB, cacheClient *redis.Client, logger log.Logger) (*Data, func(), error) {
	logHelper := log.NewHelper(log.With(logger, "module", "data/comment"))
	// 关闭Redis连接
	cleanup := func() {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := cacheClient.Ping(context.Background()).Result()
			if err != nil {
				logHelper.Warn("Redis connection pool is empty")
				return
			}
			if err = cacheClient.Close(); err != nil {
				logHelper.Errorf("Redis connection closure failed, err: %w", err)
			}
			logHelper.Info("Successfully close the Redis connection")
		}()
		wg.Wait()
	}

	data := &Data{
		db:    db,
		cache: cacheClient,
		log:   logHelper,
	}
	return data, cleanup, nil
}

// NewMysqlConn mysql数据库连接
func NewMysqlConn(c *conf.Data) *gorm.DB {
	db, err := gorm.Open(mysql.Open(c.Mysql.Dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Database connection failure, err : %v", err)
	}
	InitDB(db)
	log.Info("Database enabled successfully!")
	return db.Model(&Followers{})
}

// NewRedisConn Redis数据库连接
func NewRedisConn(c *conf.Data) *redis.Client {
	client := redis.NewClient(&redis.Options{
		DB:           int(c.Redis.FollowerDb),
		Addr:         c.Redis.Addr,
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		Password:     c.Redis.Password,
	})

	// ping Redis客户端，判断连接是否存在
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Redis database connection failure, err : %v", err)
	}
	log.Info("Cache enabled successfully!")
	return client
}

// InitDB 创建followers数据表，并自动迁移
func InitDB(db *gorm.DB) {
	if err := db.AutoMigrate(&Followers{}); err != nil {
		log.Fatalf("Database initialization error, err : %v", err)
	}
}
