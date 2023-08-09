package data

import (
	"Atreus/app/user/service/internal/conf"
	"Atreus/pkg/gorms"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGormDb, NewUserRepo)

// Data .
type Data struct {
	db  *gorms.GormConn
	log *log.Helper
}

// NewData .
func NewData(db *gorms.GormConn, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	data := &Data{
		db:  db,
		log: log.NewHelper(logger),
	}
	return data, cleanup, nil
}

// NewGormDb .
func NewGormDb(c *conf.Data) *gorms.GormConn {
	dsn := c.Database.Source
	open, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("database connect failed, error: " + err.Error())
	}
	db, _ := open.DB()
	//连接池配置
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)
	return gorms.New(open)
}
