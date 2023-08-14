package data

import (
	"Atreus/app/favorite/service/internal/conf"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sync"
)

var ProviderSet = wire.NewSet(NewData, NewFavoriteRepo, NewUserRepo, NewPublishRepo, NewMysqlConn, NewRedisConn, NewTransaction)

// CacheClient favorite 服务的 Redis 缓存客户端
type CacheClient struct {
	favoriteNumber *redis.Client // 用户点赞数缓存
	favoriteCache  *redis.Client // 用户点赞关系缓存
}

type Data struct {
	db    *gorm.DB
	cache *CacheClient
	log   *log.Helper
}

func NewData(db *gorm.DB, cacheClient *CacheClient, logger log.Logger) (*Data, func(), error) {
	logHelper := log.NewHelper(log.With(logger, "module", "data/favorite"))
	// 并发关闭所有数据库连接，后期根据Redis与Mysql是否数据同步修改
	// MySQL 会自动关闭连接，但是 Redis 不会，所以需要手动关闭
	cleanup := func() {
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			_, err := cacheClient.favoriteNumber.Ping(context.Background()).Result()
			if err != nil {
				return
			}
			if err = cacheClient.favoriteNumber.Close(); err != nil {
				logHelper.Errorf("Redis connection closure failed, err: %w", err)
			}
			logHelper.Info("Successfully close the Redis connection")
		}()
		go func() {
			defer wg.Done()
			_, err := cacheClient.favoriteCache.Ping(context.Background()).Result()
			if err != nil {
				return
			}
			if err = cacheClient.favoriteCache.Close(); err != nil {
				logHelper.Errorf("Redis connection closure failed, err: %w", err)
			}
			logHelper.Info("Successfully close the Redis connection")
		}()
		wg.Wait()
	}

	data := &Data{
		db:    db.Model(&Favorite{}), // specify table in advance
		cache: cacheClient,
		log:   logHelper,
	}
	return data, cleanup, nil
}

// NewMysqlConn mysql数据库连接
func NewMysqlConn(c *conf.Data, l log.Logger) *gorm.DB {
	logs := log.NewHelper(log.With(l, "module", "data/mysql"))
	db, err := gorm.Open(mysql.Open(c.Mysql.Dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		logs.Fatalf("Database connection failure, err : %v", err)
	}
	InitDB(db)
	logs.Info("Database enabled successfully!")
	return db.Model(&Favorite{})
}

// NewRedisConn Redis数据库连接, 并发开启连接提高速率
func NewRedisConn(c *conf.Data, l log.Logger) (cacheClient *CacheClient) {
	logs := log.NewHelper(log.With(l, "module", "data/redis"))
	// 初始化点赞数Redis客户端
	numberDB := redis.NewClient(&redis.Options{
		DB:           int(c.Redis.FavoriteNumberDb),
		Addr:         c.Redis.Addr,
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		Password:     c.Redis.Password,
	})
	// ping Redis客户端，判断连接是否存在
	_, err := numberDB.Ping(context.Background()).Result()
	if err != nil {
		logs.Fatalf("Redis database connection failure, err : %v", err)
	}
	// 初始化点赞关系Redis客户端
	relationDB := redis.NewClient(&redis.Options{
		DB:           int(c.Redis.FavoriteCacheDb),
		Addr:         c.Redis.Addr,
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		Password:     c.Redis.Password,
	})
	// ping Redis客户端，判断连接是否存在
	_, err = relationDB.Ping(context.Background()).Result()
	if err != nil {
		logs.Fatalf("Redis database connection failure, err : %v", err)
	}

	logs.Info("Cache enabled successfully!")
	return &CacheClient{favoriteNumber: numberDB, favoriteCache: relationDB}
}

// InitDB 创建User数据表，并自动迁移
func InitDB(db *gorm.DB) {
	if err := db.AutoMigrate(&Favorite{}); err != nil {
		log.Fatalf("Database initialization error, err : %v", err)
	}
}

// 用来承载事务的上下文
type contextTxKey struct{}

// Transaction 新增事务接口方法 - 来源：https://learnku.com/articles/65506
type Transaction interface {
	ExecTx(context.Context, func(ctx context.Context) error) error
}

// NewTransaction .
func NewTransaction(d *Data) Transaction {
	return d
}

// ExecTx gorm Transaction
func (d *Data) ExecTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, contextTxKey{}, tx)
		return fn(ctx)
	})
}

// DB 根据此方法来判断当前的 db 是不是使用 事务的 DB
func (d *Data) DB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(contextTxKey{}).(*gorm.DB)
	if ok {
		return tx
	}
	return d.db
}
