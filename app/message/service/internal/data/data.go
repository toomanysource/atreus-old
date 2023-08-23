package data

import (
	"Atreus/app/message/service/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/segmentio/kafka-go"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var ProviderSet = wire.NewSet(NewData, NewMessageRepo, NewMysqlConn, NewKafkaConn)

type KafkaConn struct {
	writer *kafka.Writer
	reader *kafka.Reader
}

type Data struct {
	db  *gorm.DB
	kfk *KafkaConn
	log *log.Helper
}

func NewData(db *gorm.DB, kfk *KafkaConn, logger log.Logger) (*Data, func(), error) {
	logHelper := log.NewHelper(log.With(logger, "module", "data/comment"))

	cleanup := func() {
		if err := kfk.writer.Close(); err != nil {
			logHelper.Errorf("Kafka connection closure failed, err: %w", err)
		}
		if err := kfk.reader.Close(); err != nil {
			logHelper.Errorf("Kafka connection closure failed, err: %w", err)
		}
		logHelper.Info("[kafka] client stopping")
		return
	}

	data := &Data{
		db:  db.Model(&Message{}),
		kfk: kfk,
		log: logHelper,
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
	return db
}

func NewKafkaConn(c *conf.Data) *KafkaConn {
	writer := kafka.Writer{
		Addr:                   kafka.TCP(c.Kafka.Addr),
		Topic:                  c.Kafka.Topic,
		Balancer:               &kafka.LeastBytes{},
		WriteTimeout:           c.Kafka.WriteTimeout.AsDuration(),
		ReadTimeout:            c.Kafka.ReadTimeout.AsDuration(),
		AllowAutoTopicCreation: true,
	}
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{c.Kafka.Addr},
		Partition: int(c.Kafka.Partition),
		GroupID:   "store",
		Topic:     c.Kafka.Topic,
		MaxBytes:  10e6, // 10MB
	})
	log.Info("Kafka enabled successfully!")
	return &KafkaConn{
		writer: &writer,
		reader: reader,
	}
}

// InitDB 创建followers数据表，并自动迁移
func InitDB(db *gorm.DB) {
	if err := db.AutoMigrate(&Message{}); err != nil {
		log.Fatalf("Database initialization error, err : %v", err)
	}
}
