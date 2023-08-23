package biz

import (
	"Atreus/app/message/service/internal/conf"
	"Atreus/pkg/common"
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
)

type Message struct {
	Id         uint32
	ToUserId   uint32
	FromUserId uint32
	Content    string
	CreateTime int64
}

type MessageRepo interface {
	GetMessageList(context.Context, uint32, uint32, int64) ([]*Message, error)
	PublishMessage(context.Context, uint32, uint32, string) error
	InitStoreMessageQueue()
}

type MessageUsecase struct {
	repo MessageRepo
	conf *conf.JWT
	log  *log.Helper
}

func NewMessageUsecase(repo MessageRepo, conf *conf.JWT, logger log.Logger) *MessageUsecase {
	go repo.InitStoreMessageQueue()
	return &MessageUsecase{
		repo: repo, conf: conf,
		log: log.NewHelper(log.With(logger, "model", "usecase/message")),
	}
}

func (uc *MessageUsecase) GetMessageList(
	ctx context.Context, tokenString string, toUserId uint32, preMsgTime int64) ([]*Message, error) {
	token, err := common.ParseToken(uc.conf.Http.TokenKey, tokenString)
	if err != nil {
		return nil, err
	}
	data, err := common.GetTokenData(token)
	if err != nil {
		return nil, err
	}
	userId := uint32(data["user_id"].(float64))
	return uc.repo.GetMessageList(ctx, userId, toUserId, preMsgTime)
}

func (uc *MessageUsecase) PublishMessage(ctx context.Context, tokenString string, toUserId uint32, actionType uint32, content string) error {
	token, err := common.ParseToken(uc.conf.Http.TokenKey, tokenString)
	if err != nil {
		return err
	}
	data, err := common.GetTokenData(token)
	if err != nil {
		return err
	}
	userId := uint32(data["user_id"].(float64))
	switch actionType {
	case 1:
		return uc.repo.PublishMessage(ctx, userId, toUserId, content)
	default:
		return errors.New("the actionType value for the error is provided")
	}
}
