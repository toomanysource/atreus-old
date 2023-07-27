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
type Publish struct {
	video_id int
	user_id  int
}

// GreeterRepo is a Greater repo.
type PublishRepo interface {
	Save(context.Context, *Publish) (*Publish, error)
	Update(context.Context, *Publish) (*Publish, error)
	FindByID(context.Context, int64) (*Publish, error)
	ListByHello(context.Context, string) ([]*Publish, error)
	ListAll(context.Context) ([]*Publish, error)
}

// GreeterUsecase is a Greeter usecase.
type PublishUsecase struct {
	repo PublishRepo
	log  *log.Helper
}

// NewGreeterUsecase new a Greeter usecase.
func NewPublishUsecase(repo PublishRepo, logger log.Logger) *PublishUsecase {
	return &PublishUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateGreeter creates a Greeter, and returns the new Greeter.
func (uc *PublishUsecase) CreateGreeter(ctx context.Context, g *Publish) (*Publish, error) {
	uc.log.WithContext(ctx).Infof("CreatePublish: %v", g.user_id)
	return uc.repo.Save(ctx, g)
}
