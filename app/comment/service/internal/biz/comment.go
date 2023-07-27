package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/genproto/googleapis/type/datetime"
)

// var (
// 	// ErrUserNotFound is user not found.
// 	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
// )

// Greeter is a Greeter model.
type Comment struct {
	video_id    int
	user_id     int
	content     string
	create_date datetime.DateTime
}

// GreeterRepo is a Greater repo.
type CommentRepo interface {
	Save(context.Context, *Comment) (*Comment, error)
	Update(context.Context, *Comment) (*Comment, error)
	FindByID(context.Context, int64) (*Comment, error)
	ListByHello(context.Context, string) ([]*Comment, error)
	ListAll(context.Context) ([]*Comment, error)
}

// GreeterUsecase is a Greeter usecase.
type CommentUsecase struct {
	repo CommentRepo
	log  *log.Helper
}

// NewGreeterUsecase new a Greeter usecase.
func NewCommentUsecase(repo CommentRepo, logger log.Logger) *CommentUsecase {
	return &CommentUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateGreeter creates a Greeter, and returns the new Greeter.
func (uc *CommentUsecase) CreateGreeter(ctx context.Context, g *Comment) (*Comment, error) {
	uc.log.WithContext(ctx).Infof("CreateComment: %v", g.user_id)
	return uc.repo.Save(ctx, g)
}
