package data

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"strconv"
	"time"

	"Atreus/app/comment/service/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// Comment Database Model
type Comment struct {
	Id       uint32 `gorm:"primary_key"`
	UserId   uint32 `gorm:"column:user_id;not null"`
	VideoId  uint32 `gorm:"column:video_id;not null"`
	Content  string `gorm:"column:content;not null"`
	CreateAt string `gorm:"column:created_at;default:''"`
	gorm.DeletedAt
}

func (Comment) TableName() string {
	return "comments"
}

type commentRepo struct {
	data     *Data
	userRepo *UserRepo
	log      *log.Helper
}

func NewCommentRepo(data *Data, conn *grpc.ClientConn, logger log.Logger) biz.CommentRepo {
	return &commentRepo{
		data:     data,
		userRepo: NewUserRepo(conn),
		log:      log.NewHelper(log.With(logger, "model", "comment-service/repo")),
	}
}

// DeleteComment 删除评论，先在数据库中删除，再在redis缓存中删除串行执行，保证数据一致性
func (r *commentRepo) DeleteComment(
	ctx context.Context, videoId, commentId uint32, userId uint32) (c *biz.Comment, err error) {
	c, err = r.DelComment(ctx, videoId, commentId, userId)
	if err != nil {
		return nil, err
	}
	//在redis缓存中查询评论是否存在
	comment, err := r.data.cache.HGet(
		ctx, strconv.Itoa(int(videoId)), strconv.Itoa(int(commentId))).Result()
	if err == nil {
		co := &biz.Comment{}
		err = json.Unmarshal([]byte(comment), co)
		if err != nil {
			return nil, fmt.Errorf("json unmarshal error, err : %w", err)
		}
		err = r.data.cache.HDel(ctx, strconv.Itoa(int(videoId)), strconv.Itoa(int(commentId))).Err()
		if err != nil {
			return nil, fmt.Errorf("redis transaction error, err : %w", err)
		}
	}
	return c, nil
}

// CreateComment 创建评论
func (r *commentRepo) CreateComment(
	ctx context.Context, videoId uint32, commentText string, userId uint32) (c *biz.Comment, err error) {
	// 先在数据库中插入评论
	c, err = r.InsertComment(ctx, videoId, commentText, userId)
	if err != nil {
		return nil, err
	}
	// 在redis缓存中查询是否存在视频评论列表
	err = r.data.cache.HLen(ctx, strconv.Itoa(int(videoId))).Err()
	if err != nil {
		// 如果不存在则创建
		cl, err := r.SearchCommentList(ctx, videoId)
		if err != nil {
			return nil, err
		}
		err = r.CacheCreateCommentTransaction(ctx, cl, videoId)
		if err != nil {
			return nil, err
		}
	} else {
		// 将评论存入redis缓存
		err = r.data.cache.HSet(ctx, strconv.Itoa(int(videoId)), strconv.Itoa(int(c.Id)), c).Err()
		if err != nil {
			return nil, fmt.Errorf("redis transaction error, err : %w", err)
		}
	}
	return c, nil
}

// GetCommentList 获取评论列表
func (r *commentRepo) GetCommentList(
	ctx context.Context, videoId uint32) (cl []*biz.Comment, err error) {
	// 先在redis缓存中查询是否存在视频评论列表
	commentList, err := r.data.cache.HGetAll(ctx, strconv.Itoa(int(videoId))).Result()
	if err == nil {
		// 如果存在则直接返回
		for _, v := range commentList {
			co := &biz.Comment{}
			err = json.Unmarshal([]byte(v), co)
			if err != nil {
				return nil, fmt.Errorf("json unmarshal error, err : %w", err)
			}
			cl = append(cl, co)
		}
		return cl, nil
	} else {
		// 如果不存在则创建
		cl, err := r.SearchCommentList(ctx, videoId)
		if err != nil {
			return nil, err
		}
		err = r.CacheCreateCommentTransaction(ctx, cl, videoId)
		if err != nil {
			return nil, err
		}
	}
	return cl, err
}

// GetCommentNumber 获取评论总数
func (r *commentRepo) GetCommentNumber(ctx context.Context, videoId uint32) (count int64, err error) {
	// 先在redis缓存中查询是否存在视频评论列表
	count, err = r.data.cache.HLen(ctx, strconv.Itoa(int(videoId))).Result()
	if err != nil {
		// 如果不存在则创建
		cl, err := r.SearchCommentList(ctx, videoId)
		if err != nil {
			return 0, err
		}
		err = r.CacheCreateCommentTransaction(ctx, cl, videoId)
		if err != nil {
			return 0, err
		}
	}
	return count, nil
}

// DelComment 数据库删除评论
func (r *commentRepo) DelComment(
	ctx context.Context, videoId, commentId uint32, userId uint32) (c *biz.Comment, err error) {

	comment := &Comment{}
	result := r.data.db.WithContext(ctx).First(comment, commentId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("comment don't exist")
	} else if result.Error != nil {
		return nil, fmt.Errorf("query error, err : %w", result.Error)
	}
	// 判断当前用户是否为评论用户
	if comment.UserId != userId {
		return nil, errors.New("mismatch between commenter id and user id")
	}
	// 判断视频id是否为当前视频id
	if comment.VideoId != videoId {
		return nil, errors.New("comment video id doesn't match current video id")
	}
	result = r.data.db.WithContext(ctx).Delete(&Comment{}, commentId)
	if err = result.Error; err != nil {
		return nil, fmt.Errorf("an error occurs when deleting, err : %w", err)
	}
	return nil, nil
}

// InsertComment 数据库插入评论
func (r *commentRepo) InsertComment(
	ctx context.Context, videoId uint32, commentText string, userId uint32) (*biz.Comment, error) {

	if commentText == "" {
		return nil, errors.New("content are empty")
	}

	users, err := r.userRepo.GetUserInfoByUserIds(ctx, []uint32{userId})
	if err != nil {
		return nil, fmt.Errorf("user service transfer error, err : %w", err)
	}

	comment := &Comment{
		UserId:   userId,
		VideoId:  videoId,
		Content:  commentText,
		CreateAt: time.Now().Format("01-02"),
	}

	result := r.data.db.WithContext(ctx).Create(comment)
	if err := result.Error; err != nil {
		return nil, fmt.Errorf("an error occurred while creating the comment, err : %w", err)
	}

	return &biz.Comment{
		Id: comment.Id,
		User: &biz.User{
			Id:              users[0].Id,
			Name:            users[0].Name,
			Avatar:          users[0].Avatar,
			BackgroundImage: users[0].BackgroundImage,
			Signature:       users[0].Signature,
			IsFollow:        false,
			FollowCount:     users[0].FollowCount,
			FollowerCount:   users[0].FollowerCount,
			TotalFavorited:  users[0].TotalFavorited,
			WorkCount:       users[0].WorkCount,
			FavoriteCount:   users[0].FavoriteCount,
		},
		Content:    commentText,
		CreateDate: comment.CreateAt,
	}, nil
}

// SearchCommentList 数据库搜索评论列表
func (r *commentRepo) SearchCommentList(
	ctx context.Context, videoId uint32) (cl []*biz.Comment, err error) {

	var commentList []*Comment
	result := r.data.db.WithContext(ctx).Where("video_id = ?", videoId).
		Order(gorm.Expr("STR_TO_DATE(create_at, '%m-%d') DESC")).Find(commentList)
	if err := result.Error; err != nil {
		return nil, fmt.Errorf("an error occurs when the query, err : %w", err)
	}

	// 此视频没有评论
	if result.RowsAffected == 0 {
		return nil, nil
	}

	// 获取评论列表中的所有用户id
	userIds := make([]uint32, 0, len(commentList)+1)
	for _, comment := range commentList {
		userIds = append(userIds, comment.UserId)
	}

	// 统一查询，减少网络IO
	users, err := r.userRepo.GetUserInfoByUserIds(ctx, userIds)
	if err != nil {
		return nil, err
	}

	// 返回的数据可能乱序，映射map
	userMap := make(map[uint32]*biz.User)
	for _, user := range users {
		userMap[user.Id] = user
	}

	for _, comment := range commentList {
		cl = append(cl, &biz.Comment{
			Id:         comment.Id,
			User:       userMap[comment.UserId],
			Content:    comment.Content,
			CreateDate: comment.CreateAt,
		})
	}
	return cl, nil
}

// CountCommentNumber 数据库统计视频总评论数
func (r *commentRepo) CountCommentNumber(ctx context.Context, videoId uint32) (count int64, err error) {
	result := r.data.db.WithContext(ctx).Where("video_id = ?", videoId).Count(&count)
	if err = result.Error; err != nil {
		return 0, fmt.Errorf("error in counting quantity, err: %w", err)
	}
	return count, err
}
