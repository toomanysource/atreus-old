package data

import (
	"Atreus/app/comment/service/internal/biz"
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func (r *commentRepo) CacheCreateCommentTransaction(ctx context.Context, cl []*biz.Comment, videoId uint32) error {
	// 使用事务将评论列表存入redis缓存
	pipe := r.data.cache.Pipeline()
	insertMap := make(map[string]interface{}, len(cl))
	for _, v := range cl {
		insertMap[strconv.Itoa(int(v.Id))] = v
	}
	pipe.HMSet(ctx, strconv.Itoa(int(videoId)), insertMap)
	// 将评论数量存入redis缓存,使用随机过期时间防止缓存雪崩
	random := rand.Intn(361) + 360
	pipe.Expire(ctx, strconv.Itoa(int(videoId)), time.Minute*time.Duration(random))
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("redis transaction error, err : %w", err)
	}
	return err
}
