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
type Relation struct {
	video_id int
	user_id  int
}

// GreeterRepo is a Greater repo.
type RelationRepo interface {
	Save(context.Context, *Relation) (*Relation, error)
	Update(context.Context, *Relation) (*Relation, error)
	FindByID(context.Context, int64) (*Relation, error)
	ListByHello(context.Context, string) ([]*Relation, error)
	ListAll(context.Context) ([]*Relation, error)
}

// GreeterUsecase is a Greeter usecase.
type RelationUsecase struct {
	repo RelationRepo
	log  *log.Helper
}

// NewGreeterUsecase new a Greeter usecase.
func NewRelationUsecase(repo RelationRepo, logger log.Logger) *RelationUsecase {
	return &RelationUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateGreeter creates a Greeter, and returns the new Greeter.
func (uc *RelationUsecase) CreateGreeter(ctx context.Context, g *Relation) (*Relation, error) {
	uc.log.WithContext(ctx).Infof("CreateRelation: %v", g.user_id)
	return uc.repo.Save(ctx, g)
}
