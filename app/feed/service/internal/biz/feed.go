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
type Feed struct {
	video_id int
	user_id  int
}

// GreeterRepo is a Greater repo.
type FeedRepo interface {
	Save(context.Context, *Feed) (*Feed, error)
	Update(context.Context, *Feed) (*Feed, error)
	FindByID(context.Context, int64) (*Feed, error)
	ListByHello(context.Context, string) ([]*Feed, error)
	ListAll(context.Context) ([]*Feed, error)
}

// GreeterUsecase is a Greeter usecase.
type FeedUsecase struct {
	repo FeedRepo
	log  *log.Helper
}

// NewGreeterUsecase new a Greeter usecase.
func NewFeedUsecase(repo FeedRepo, logger log.Logger) *FeedUsecase {
	return &FeedUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateGreeter creates a Greeter, and returns the new Greeter.
func (uc *FeedUsecase) CreateGreeter(ctx context.Context, g *Feed) (*Feed, error) {
	uc.log.WithContext(ctx).Infof("CreateFeed: %v", g.user_id)
	return uc.repo.Save(ctx, g)
}
