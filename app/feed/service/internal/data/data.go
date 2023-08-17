package data

import (
	"Atreus/app/feed/service/internal/conf"
	"context"
	"sync"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var ProviderSet = wire.NewSet(NewData, NewFeeedRepo, NewDB, NewRedisConn)

// Kafka
// type message struct {
// 	writer *kafka.Writer
// 	reader *kafka.Reader
// }

type Data struct {
	db    *gorm.DB
	cache *redis.Client
	// mq    *message
	log *log.Helper
}

func NewData(db *gorm.DB, feedCacheClient *redis.Client, logger log.Logger) (*Data, func(), error) {
	logHelper := log.NewHelper(log.With(logger, "module", "data/feed"))

	// concurrent close all DB connect.
	cleanup := func() {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := feedCacheClient.Ping(context.Background()).Result()
			if err != nil {
				logHelper.Warn("Redis connect poll is empty.")
				return
			}
			if err = feedCacheClient.Close(); err != nil {
				logHelper.Errorf("Redis connect closure failed, err: %w", err)
			}
			logHelper.Info("Success close the Redis connect.")
		}()
		wg.Wait()
	}
	data := &Data{
		db:    db,
		cache: feedCacheClient,
		log:   logHelper,
	}
	return data, cleanup, nil
}

// mysql Database Connect.
func NewDB(c *conf.Data) *gorm.DB {
	db, err := gorm.Open(mysql.Open(c.Mysql.Dsn))
	if err != nil {
		log.Fatalf("Failed to connect to database: err: %v", err)
	}
	InitDB(db)
	log.Info("Database is running successfully.")
	return db
}

// Create Redis connect. Tests that a connect exists.
func NewRedisConn(c *conf.Data) *redis.Client {
	client := redis.NewClient(&redis.Options{
		DB:           int(c.Redis.Db), // tbd.
		Addr:         c.Redis.Addr,
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		Password:     c.Redis.Password,
	})

	// ping Redis client. Test
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Redis database connection failure, err : %v", err)
	}
	log.Info("CommentNumberCache enabled successfully!")
	return client
}

// InitDB, create the correspond table and auto migrate.
func InitDB(db *gorm.DB) {
	if err := db.AutoMigrate(&Video{}); err != nil {
		log.Fatalf("Database init error, err: %v", err)
	}
}
