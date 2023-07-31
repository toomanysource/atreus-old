package data

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"time"

	"Atreus/app/comment/service/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

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

// DeleteComment 删除评论
func (r *commentRepo) DeleteComment(
	ctx context.Context, videoId, commentId uint32, userId uint32) (c *biz.Comment, err error) {
	comment := &Comment{}
	result := r.data.db.WithContext(ctx).First(comment, commentId)
	if err = result.Error; err != nil {
		return nil, err
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
		return nil, err
	}
	return nil, nil
}

// CreateComment 创建评论
func (r *commentRepo) CreateComment(
	ctx context.Context, videoId uint32, commentText string, userId uint32) (*biz.Comment, error) {

	if commentText == "" {
		return nil, errors.New("content are empty")
	}
	users, err := r.userRepo.GetUserInfoByUserId(ctx, []uint32{userId})
	if err != nil {
		return nil, err
	}

	comment := &Comment{
		UserId:   userId,
		VideoId:  videoId,
		Content:  commentText,
		CreateAt: time.Now().Format("01-02"),
	}

	result := r.data.db.WithContext(ctx).Create(comment)
	if err := result.Error; err != nil {
		return nil, err
	}

	// jwt-go解析的payload,整型数据被定义为float64类型,因此需要先断言为float64,再强转uint32
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

// GetCommentList 获取评论列表
func (r *commentRepo) GetCommentList(
	ctx context.Context, videoId uint32) (cl []*biz.Comment, err error) {

	var commentList []*Comment
	r.data.db.WithContext(ctx).Where("video_id = ?", videoId).Find(commentList)

	// 获取评论列表中的所有用户id
	userIds := make([]uint32, 0, len(commentList)+1)
	for _, comment := range commentList {
		userIds = append(userIds, comment.UserId)
	}

	// 统一查询，减少网络IO
	users, err := r.userRepo.GetUserInfoByUserId(ctx, userIds)
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