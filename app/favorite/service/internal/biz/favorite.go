package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

// var (
// 	// ErrUserNotFound is user not found.
// 	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
// )

// Greeter is a Greeter model.
type Favorite struct {
	video_id int
	user_id  int
}

// GreeterRepo is a Greater repo.
type FavoriteRepo interface {
	Save(context.Context, *favorite) (*favorite, error)
	Update(context.Context, *favorite) (*favorite, error)
	FindByID(context.Context, int64) (*favorite, error)
	ListByHello(context.Context, string) ([]*favorite, error)
	ListAll(context.Context) ([]*favorite, error)
}

// GreeterUsecase is a Greeter usecase.
type FavoriteUsecase struct {
	repo FavoriteRepo
	log  *log.Helper
}

// NewGreeterUsecase new a Greeter usecase.
func NewGreeterUsecase(repo FavoriteRepo, logger log.Logger) *FavoriteUsecase {
	return &FavoriteUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateGreeter creates a Greeter, and returns the new Greeter.
func (uc *FavoriteUsecase) CreateGreeter(ctx context.Context, g *Greeter) (*Greeter, error) {
	uc.log.WithContext(ctx).Infof("CreateGreeter: %v", g.Hello)
	return uc.repo.Save(ctx, g)
}
