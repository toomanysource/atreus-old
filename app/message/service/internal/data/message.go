package data

import (
	"Atreus/app/message/service/internal/biz"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/segmentio/kafka-go"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type Message struct {
	Id         uint32 `gorm:"column:id;primary_key;auto_increment"`
	FromUserId uint32 `gorm:"column:from_user_id;not null"`
	ToUserId   uint32 `gorm:"column:to_user_id;not null"`
	Content    string `gorm:"column:content;not null"`
	CreateAt   int64  `gorm:"column:created_at"`
}

func (Message) TableName() string {
	return "message"
}

type messageRepo struct {
	data *Data
	log  *log.Helper
}

func NewMessageRepo(data *Data, logger log.Logger) biz.MessageRepo {
	return &messageRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *messageRepo) GetMessageList(ctx context.Context, userId uint32, toUSerId uint32, preMsgTime int64) ([]*biz.Message, error) {
	return r.SearchMessage(ctx, userId, toUSerId, preMsgTime)
}

// PublishMessage 生产者发送消息
func (r *messageRepo) PublishMessage(ctx context.Context, userId, toUSerId uint32, content string) error {
	if userId == toUSerId {
		return errors.New("can't send message to yourself")
	}
	mg := &Message{
		FromUserId: userId,
		ToUserId:   toUSerId,
		Content:    content,
		CreateAt:   time.Now().UnixMilli(),
	}
	byteValue, err := json.Marshal(mg)
	if err != nil {
		return fmt.Errorf("json marshal error, err: %w", err)
	}
	err = r.data.kfk.writer.WriteMessages(context.TODO(), kafka.Message{
		Partition: 0,
		Value:     byteValue,
	})
	if err != nil {
		return fmt.Errorf("write message error, err: %w", err)
	}
	return nil
}

// SearchMessage 数据库根据最新消息时间查询消息
func (r *messageRepo) SearchMessage(ctx context.Context, userId, toUSerId uint32, preMsgTime int64) (ml []*biz.Message, err error) {
	var mel []*Message
	err = r.data.db.WithContext(ctx).Where(
		"(from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?)",
		userId, toUSerId, toUSerId, userId).Where("created_at > ?", preMsgTime).
		Order("created_at").Find(&mel).Error
	if err != nil {
		return nil, fmt.Errorf("search message error, err: %w", err)
	}
	if err = copier.Copy(&ml, &mel); err != nil {
		return nil, fmt.Errorf("copy message error, err: %w", err)
	}
	return
}

// InitStoreMessageQueue 初始化聊天记录存储队列
func (r *messageRepo) InitStoreMessageQueue() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// 监听Ctrl+C退出信号
	signChan := make(chan os.Signal, 1)
	signal.Notify(signChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signChan
		cancel()
	}()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := r.data.kfk.reader.ReadMessage(ctx)
			if errors.Is(err, context.Canceled) {
				return
			}
			if err != nil {
				r.log.Errorf("read message error, err: %v", err)
			}
			value := msg.Value
			var mg *Message
			err = json.Unmarshal(value, &mg)
			if err != nil {
				r.log.Errorf("json unmarshal error, err: %v", err)
				return
			}
			err = r.InsertMessage(ctx, mg.FromUserId, mg.ToUserId, mg.Content)
			if err != nil {
				r.log.Errorf("insert message error, err: %v", err)
				return
			}
			err = r.data.kfk.reader.CommitMessages(ctx, msg)
			if err != nil {
				r.log.Errorf("commit message error, err: %v", err)
				return
			}
			r.log.Infof("message: UserId-%v to UserId-%v: %v ", mg.FromUserId, mg.ToUserId, mg.Content)
		}
	}
}

// InsertMessage 数据库插入消息
func (r *messageRepo) InsertMessage(ctx context.Context, userId uint32, toUSerId uint32, content string) error {
	err := r.data.db.WithContext(ctx).Create(&Message{
		FromUserId: userId,
		ToUserId:   toUSerId,
		Content:    content,
		CreateAt:   time.Now().UnixMilli(),
	}).Error
	if err != nil {
		return fmt.Errorf("insert message error, err: %w", err)
	}
	return nil
}
