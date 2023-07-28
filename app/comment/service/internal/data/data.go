package data

import (
	"Atreus/app/comment/service/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var ProviderSet = wire.NewSet(NewData, NewCommentRepo, NewMysqlConn)

type Data struct {
	db  *gorm.DB
	log *log.Helper
}

// NewMysqlConn mysql数据库连接
func NewMysqlConn(c *conf.Data) *gorm.DB {
	db, err := gorm.Open(mysql.Open(c.Mysql.Dsn))
	if err != nil {
		panic(err)
	}
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
			logHelper.Error("Mysql connection closure failed, err: ", err)
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
