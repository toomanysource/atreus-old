package biz

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-jwt/jwt/v4"
)

type Comment struct {
	Id         uint32
	User       *User
	Content    string
	CreateDate string
}

type User struct {
	Id              uint32
	Name            string
	Avatar          string
	BackgroundImage string
	Signature       string
	IsFollow        bool
	FollowCount     uint32
	FollowerCount   uint32
	TotalFavorited  uint32
	WorkCount       uint32
	FavoriteCount   uint32
}

// parseToken 接收TokenString进行校验
func (uc *CommentUsecase) parseToken(tokenKey, tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenKey), nil
	})
	if err != nil {
		log.Errorf("Server failed to convert Token, err :", err.Error())
		return nil, err
	}
	if token.Valid {
		return token, nil
	}
	return nil, errors.New("invalid JWT token")
}

// getTokenData 获取Token中的用户数据,返回的是map[string]any类型，需要断言
func (uc *CommentUsecase) getTokenData(token *jwt.Token) (map[string]any, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok {
		return claims, nil
	}
	return nil, errors.New("failed to extract claims from JWT token")
}

func (uc *CommentUsecase) GetCommentList(
	ctx context.Context, tokenString string, videoId uint32) ([]*Comment, error) {
	token, err := uc.parseToken(uc.config.Http.TokenKey, tokenString)
	if err != nil {
		return nil, err
	}
	_, err = uc.getTokenData(token)
	if err != nil {
		return nil, err
	}
	return uc.commentRepo.GetCommentList(ctx, videoId)
}

func (uc *CommentUsecase) CommentAction(
	ctx context.Context, videoId, commentId uint32,
	actionType uint32, commentText string, tokenString string) (*Comment, error) {
	token, err := uc.parseToken(uc.config.Http.TokenKey, tokenString)
	if err != nil {
		return nil, err
	}
	data, err := uc.getTokenData(token)
	if err != nil {
		return nil, err
	}
	userId := uint32(data["user_id"].(float64))
	switch actionType {
	case 1:
		return uc.commentRepo.CreateComment(ctx, videoId, commentText, userId)
	case 2:
		return uc.commentRepo.DeleteComment(ctx, videoId, commentId, userId)
	default:
		return nil, errors.New("the value of action_type is not in the specified range")
	}
}

func (uc *CommentUsecase) GetCommentNumber(ctx context.Context, videoId uint32) (int64, error) {
	return uc.commentRepo.GetCommentNumber(ctx, videoId)
}
