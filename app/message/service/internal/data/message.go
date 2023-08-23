package data

import (
	"Atreus/app/message/service/internal/biz"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/segmentio/kafka-go"
	"math/rand"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"sync"
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

func (r *messageRepo) GetMessageList(ctx context.Context, userId uint32, toUserId uint32, preMsgTime int64) ([]*biz.Message, error) {
	// 先在redis缓存中查询是否存在聊天记录列表
	key := setKey(userId, toUserId)
	msgList, err := r.data.cache.HKeys(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("redis query error %w", err)
	}

	cl := make([]*biz.Message, len(msgList))
	if len(msgList) > 0 {
		// 如果存在则直接返回
		var wg sync.WaitGroup
		var mutex sync.Mutex
		errChan := make(chan error)
		for _, v := range msgList {
			wg.Add(1)
			go func(message string) {
				defer wg.Done()
				co := &biz.Message{}
				if err = json.Unmarshal([]byte(message), co); err != nil {
					errChan <- fmt.Errorf("json unmarshal error %w", err)
					return
				}
				mutex.Lock()
				cl = append(cl, co)
				mutex.Unlock()
			}(v)
		}
		wg.Wait()
		if err = <-errChan; err != nil {
			return nil, err
		}
	} else {
		cl, err = r.SearchMessage(userId, toUserId, preMsgTime)
		if err != nil {
			return nil, err
		}
		// 没有列表则不创建
		if len(cl) == 0 {
			return nil, nil
		}
		go func(l []*biz.Message) {
			if err = r.CacheCreateMessageTransaction(context.Background(), l, userId, toUserId); err != nil {
				r.log.Errorf("redis transaction error %w", err)
				return
			}
			r.log.Info("redis transaction success")
		}(cl)
	}
	sortMessage(cl)
	return cl, nil
}

// PublishMessage 发送消息
func (r *messageRepo) PublishMessage(ctx context.Context, userId, toUserId uint32, content string) error {
	if userId == toUserId {
		return errors.New("can't send message to yourself")
	}
	createTime := time.Now().UnixMilli()
	err := r.MessageProducer(userId, toUserId, content, createTime)
	if err != nil {
		return fmt.Errorf("message producer error, err: %w", err)
	}
	go func() {
		ctx = context.TODO()
		// 在redis缓存中查询是否存在
		key := setKey(userId, toUserId)
		count, err := r.data.cache.Exists(ctx, key).Result()
		if err != nil {
			r.log.Errorf("redis query error %w", err)
			return
		}
		if count == 0 {
			// 如果不存在则创建
			ml, err := r.SearchMessage(userId, toUserId, createTime)
			if err != nil {
				r.log.Errorf("mysql query error %w", err)
				return
			}
			// 没有聊天列表则不创建
			if len(ml) == 0 {
				return
			}
			if err = r.CacheCreateMessageTransaction(ctx, ml, userId, toUserId); err != nil {
				r.log.Errorf("redis transaction error %w", err)
				return
			}
			r.log.Info("redis transaction success")
			return
		} else {
			mg, err := r.SearchMessage(userId, toUserId, createTime)
			if err != nil {
				r.log.Errorf("mysql query error %w", err)
				return
			}
			data, err := json.Marshal(mg)
			if err != nil {
				r.log.Errorf("json marshal error %w", err)
				return
			}
			if err = r.data.cache.HSet(
				ctx, key, string(data), "").Err(); err != nil {
				r.log.Errorf("redis store error %w", err)
				return
			}
			r.log.Info("redis store success")
		}
	}()
	return nil
}

// SearchMessage 数据库根据最新消息时间查询消息
func (r *messageRepo) SearchMessage(userId, toUserId uint32, preMsgTime int64) (ml []*biz.Message, err error) {
	var mel []*Message
	err = r.data.db.Where(
		"(from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?)",
		userId, toUserId, toUserId, userId).Where("created_at > ?", preMsgTime).
		Order("created_at").Find(&mel).Error
	if err != nil {
		return nil, fmt.Errorf("search message error, err: %w", err)
	}
	for _, v := range mel {
		ml = append(ml, &biz.Message{
			Id:         v.Id,
			FromUserId: v.FromUserId,
			ToUserId:   v.ToUserId,
			Content:    v.Content,
			CreateTime: v.CreateAt,
		})
	}
	return
}

// MessageProducer 生产消息
func (r *messageRepo) MessageProducer(userId, toUserId uint32, content string, time int64) error {
	mg := &Message{
		FromUserId: userId,
		ToUserId:   toUserId,
		Content:    content,
		CreateAt:   time,
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
			err = r.InsertMessage(mg.FromUserId, mg.ToUserId, mg.Content)
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
func (r *messageRepo) InsertMessage(userId uint32, toUserId uint32, content string) error {
	err := r.data.db.Create(&Message{
		FromUserId: userId,
		ToUserId:   toUserId,
		Content:    content,
		CreateAt:   time.Now().UnixMilli(),
	}).Error
	if err != nil {
		return fmt.Errorf("insert message error, err: %w", err)
	}
	return nil
}

// CacheCreateMessageTransaction 缓存创建事务
func (r *messageRepo) CacheCreateMessageTransaction(ctx context.Context, ml []*biz.Message, userId, toUserId uint32) error {
	// 使用事务将列表存入redis缓存
	_, err := r.data.cache.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		insertMap := make(map[string]interface{}, len(ml))
		key := setKey(userId, toUserId)
		for _, u := range ml {
			data, err := json.Marshal(u)
			if err != nil {
				return fmt.Errorf("json marshal error, err: %w", err)
			}
			insertMap[string(data)] = ""
		}
		err := pipe.HMSet(ctx, key, insertMap).Err()
		if err != nil {
			return fmt.Errorf("redis store error, err : %w", err)
		}
		// 将评论数量存入redis缓存,使用随机过期时间防止缓存雪崩
		err = pipe.Expire(ctx, key, randomTime(time.Minute, 360, 720)).Err()
		if err != nil {
			return fmt.Errorf("redis expire error, err : %w", err)
		}
		return nil
	})
	return err
}

// randomTime 随机生成时间
func randomTime(timeType time.Duration, begin, end int) time.Duration {
	return timeType * time.Duration(rand.Intn(end-begin+1)+begin)
}

func setKey(userId, toUserId uint32) string {
	if userId > toUserId {
		userId, toUserId = toUserId, userId
	}
	return fmt.Sprint(strconv.Itoa(int(userId)), "-", strconv.Itoa(int(toUserId)))
}

// sortMessage 对聊天记录进行排序
func sortMessage(cl []*biz.Message) {
	// 对原始切片进行排序
	sort.Slice(cl, func(i, j int) bool {
		if cl[i].CreateTime <= cl[j].CreateTime {
			return true
		}
		return false
	})
}
