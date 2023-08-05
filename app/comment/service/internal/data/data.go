package data

import (
	"Atreus/app/comment/service/internal/conf"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
)

var ProviderSet = wire.NewSet(NewData, NewCommentRepo, NewUserRepo, NewMysqlConn, NewRedisConn)

type Data struct {
	db    *gorm.DB
	cache *redis.Client
	log   *log.Helper
}

func NewData(db *gorm.DB, cacheClient *redis.Client, logger log.Logger) (*Data, func(), error) {
	logHelper := log.NewHelper(log.With(logger, "module", "data/comment"))
	// 并发关闭所有数据库连接，后期根据Redis与Mysql是否数据同步修改
	cleanup := func() {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
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
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := cacheClient.Ping(context.Background()).Result()
			if err != nil {
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
	db, err := gorm.Open(mysql.Open(c.Mysql.Dsn))
	if err != nil {
		log.Fatalf("Database connection failure, err : %w", err)
	}
	InitDB(db)
	log.Info("Database enabled successfully!")
	return db
}

// NewRedisConn Redis数据库连接
func NewRedisConn(c *conf.Data) *redis.Client {
	client := redis.NewClient(&redis.Options{
		DB:           int(c.Redis.CommentDb),
		Addr:         c.Redis.Addr,
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		Password:     c.Redis.Password,
	})

	// ping Redis客户端，判断连接是否存在
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Redis database connection failure, err : %w", err)
	}
	log.Info("CommentNumberCache enabled successfully!")
	return client
}

// InitDB 创建User数据表，并自动迁移
func InitDB(db *gorm.DB) {
	if err := db.AutoMigrate(&Comment{}); err != nil {
		log.Fatalf("Database initialization error, err : %w", err)
	}
}
