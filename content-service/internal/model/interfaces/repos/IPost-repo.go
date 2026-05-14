package repos

import (
	"context"

	"github.com/CoffeeSi/social-network-microservices/content-service/internal/model"
)

type PostRepository interface {
	CreatePost(ctx context.Context, post model.Post) (model.Post, error)
	GetPost(ctx context.Context, id string) (model.Post, error)
	DeletePost(ctx context.Context, id string) error
	UpdatePost(ctx context.Context, id string, content string, mediaURLs []string) (model.Post, error)
	GetFeed(ctx context.Context, pageSize, page int32) ([]model.Post, int32, error)
	GetUserPosts(ctx context.Context, authorID string, pageSize, page int32) ([]model.Post, int32, error)
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
	IncrementCommentCount(ctx context.Context, postID string, amount int32) error
	IncrementLikeCount(ctx context.Context, postID string, amount int32) error
}
