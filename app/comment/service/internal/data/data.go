package data

import (
	"Atreus/app/comment/service/internal/conf"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var ProviderSet = wire.NewSet(NewData, NewCommentRepo, NewUserRepo, NewMysqlConn, NewRedisConn)

type Data struct {
	db    *gorm.DB
	cache *redis.Client
	log   *log.Helper
}

// NewMysqlConn mysql数据库连接
func NewMysqlConn(c *conf.Data) *gorm.DB {
	db, err := gorm.Open(mysql.Open(c.Mysql.Dsn))
	if err != nil {
		log.Fatalf("Database connection failure, err : %w", err)
	}
	InitDB(db)
	return db
}

// NewRedisConn Redis数据库连接
func NewRedisConn(c *conf.Data) *redis.Client {
	cache := redis.NewClient(&redis.Options{
		DB:           int(c.Redis.Db),
		Addr:         c.Redis.Addr,
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		Password:     c.Redis.Password,
	})

	_, err := cache.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Redis database connection failure, err : %w", err)
	}
	return cache
}

func NewData(db *gorm.DB, cache *redis.Client, logger log.Logger) (*Data, func(), error) {
	logHelper := log.NewHelper(log.With(logger, "module", "data/comment"))

	cleanup := func() {
		sqlDB, err := db.DB()
		// 如果err不为空，则连接池中没有连接
		if err != nil {
			return
		}
		if err = sqlDB.Close(); err != nil {
			logHelper.Errorf("Mysql connection closure failed, err: %w", err)
			return
		}
		logHelper.Info("Successfully close the Mysql connection")
		_, err = cache.Ping(context.Background()).Result()
		if err != nil {
			return
		}
		if err = cache.Close(); err != nil {
			logHelper.Errorf("Redis connection closure failed, err: %w", err)
		}
		logHelper.Info("Successfully close the Redis connection")
	}

	data := &Data{
		db:    db,
		cache: cache,
		log:   logHelper,
	}
	return data, cleanup, nil
}

// InitDB 创建User数据表，并自动迁移
func InitDB(db *gorm.DB) {
	if err := db.AutoMigrate(&Comment{}); err != nil {
		log.Fatalf("Database initialization error, err : %w", err)
	}
}
