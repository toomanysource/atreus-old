package biz

import (
	"Atreus/app/relation/service/internal/conf"
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-jwt/jwt/v4"
)

type User struct {
	Id              uint32 // 用户id
	Name            string // 用户名称
	FollowCount     uint32 // 关注总数
	FollowerCount   uint32 // 粉丝总数
	IsFollow        bool   // true-已关注，false-未关注
	Avatar          string //用户头像
	BackgroundImage string //用户个人页顶部大图
	Signature       string //个人简介
	TotalFavorite   uint32 //获赞数量
	WorkCount       uint32 //作品数量
	FavoriteCount   uint32 //点赞数量
}

type RelationRepo interface {
	GetFollowList(context.Context, uint32) ([]*User, error)
	GetFollowerList(context.Context, uint32) ([]*User, error)
	Follow(context.Context, uint32, uint32) error
	UnFollow(context.Context, uint32, uint32) error
}

type RelationUsecase struct {
	repo   RelationRepo
	config *conf.JWT
	log    *log.Helper
}

func NewRelationUsecase(repo RelationRepo, logger log.Logger) *RelationUsecase {
	return &RelationUsecase{repo: repo, log: log.NewHelper(logger)}
}

// parseToken 接收TokenString进行校验
func (uc *RelationUsecase) parseToken(tokenKey, tokenString string) (*jwt.Token, error) {
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
func (uc *RelationUsecase) getTokenData(token *jwt.Token) (map[string]any, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok {
		return claims, nil
	}
	return nil, errors.New("failed to extract claims from JWT token")
}

// GetFollowList 获取关注列表
func (uc *RelationUsecase) GetFollowList(ctx context.Context, id uint32, tokenString string) ([]*User, error) {
	token, err := uc.parseToken(uc.config.Http.TokenKey, tokenString)
	if err != nil {
		return nil, err
	}
	_, err = uc.getTokenData(token)
	if err != nil {
		return nil, err
	}
	return uc.repo.GetFollowList(ctx, id)
}

// GetFollowerList 获取粉丝列表
func (uc *RelationUsecase) GetFollowerList(ctx context.Context, id uint32, tokenString string) ([]*User, error) {
	token, err := uc.parseToken(uc.config.Http.TokenKey, tokenString)
	if err != nil {
		return nil, err
	}
	_, err = uc.getTokenData(token)
	if err != nil {
		return nil, err
	}
	return uc.repo.GetFollowerList(ctx, id)
}

// Action 关注和取消关注
func (uc *RelationUsecase) Action(ctx context.Context, tokenString string, id uint32, actionType uint32) error {
	token, err := uc.parseToken(uc.config.Http.TokenKey, tokenString)
	if err != nil {
		return errors.New("invalid JWT token")
	}
	data, err := uc.getTokenData(token)
	if err != nil {
		return errors.New("something wrong")
	}
	userId := uint32(data["user_id"].(float64))
	switch actionType {
	//1为关注
	case 1:
		err := uc.repo.Follow(ctx, userId, id)
		if err != nil {
			return errors.New("something wrong")
		}
	//2为取消关注
	case 2:
		err := uc.repo.UnFollow(ctx, userId, id)
		if err != nil {
			return errors.New("something wrong")
		}
	}
	return nil
}
