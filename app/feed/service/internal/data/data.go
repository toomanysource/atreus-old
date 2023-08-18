package data

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewFeedRepo)

//type Data struct {
//	//db    *gorm.DB
//	//cache *redis.Client
//	log *log.Helper
//}
//
//func NewData(logger log.Logger) (*Data, func(), error) {
//	logHelper := log.NewHelper(log.With(logger, "module", "data/feed"))
//
//	// concurrent close all DB connect.
//	cleanup := func() {
//		//var wg sync.WaitGroup
//		//wg.Add(1)
//		//go func() {
//		//	defer wg.Done()
//		//	_, err := feedCacheClient.Ping(context.Background()).Result()
//		//	if err != nil {
//		//		logHelper.Warn("Redis connect poll is empty.")
//		//		return
//		//	}
//		//	if err = feedCacheClient.Close(); err != nil {
//		//		logHelper.Errorf("Redis connect closure failed, err: %w", err)
//		//	}
//		//	logHelper.Info("Success close the Redis connect.")
//		//}()
//		//wg.Wait()
//	}
//	data := &Data{
//		//db:    db,
//		//cache: feedCacheClient,
//		log: logHelper,
//	}
//	return data, cleanup, nil
//}

//// NewRedisConn Create Redis connect. Tests that a connect exists.
//func NewRedisConn(c *conf.Data) *redis.Client {
//	client := redis.NewClient(&redis.Options{
//		DB:           int(c.Redis.Db), // tbd.
//		Addr:         c.Redis.Addr,
//		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
//		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
//		Password:     c.Redis.Password,
//	})
//
//	// ping Redis client. Test
//	_, err := client.Ping(context.Background()).Result()
//	if err != nil {
//		log.Fatalf("Redis database connection failure, err : %v", err)
//	}
//	log.Info("CommentNumberCache enabled successfully!")
//	return client
//}
