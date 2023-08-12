package biz

import (
	"context"
	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewFavoriteUsecase)

// Transaction 新增事务接口方法 - 来源：https://learnku.com/articles/65506
type Transaction interface {
	ExecTx(context.Context, func(ctx context.Context) error) error
}
