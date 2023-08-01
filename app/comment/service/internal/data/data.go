package data

import (
	"Atreus/app/comment/service/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var ProviderSet = wire.NewSet(NewData, NewCommentRepo, NewUserRepo, NewMysqlConn)

type Data struct {
	db  *gorm.DB
	log *log.Helper
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

func NewData(db *gorm.DB, logger log.Logger) (*Data, func(), error) {
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
	}

	data := &Data{
		db:  db,
		log: logHelper,
	}
	return data, cleanup, nil
}

// InitDB 创建User数据表，并自动迁移
func InitDB(db *gorm.DB) {
	if err := db.AutoMigrate(&Comment{}); err != nil {
		log.Fatalf("Database initialization error, err : %w", err)
	}
}
