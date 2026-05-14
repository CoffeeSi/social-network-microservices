package repos

import (
	"context"

	"github.com/CoffeeSi/social-network-microservices/content-service/internal/model"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, comment model.Comment) (model.Comment, error)
	GetComment(ctx context.Context, id string) (model.Comment, error)
	ListComments(ctx context.Context, pID string, pageSize, page int32) ([]model.Comment, int32, error)
	UpdateComment(ctx context.Context, id string, text string) (model.Comment, error)
	DeleteComment(ctx context.Context, id string) error
}
