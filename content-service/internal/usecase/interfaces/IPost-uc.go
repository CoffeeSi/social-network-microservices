package interfaces

import (
	"context"

	"github.com/CoffeeSi/social-network-microservices/content-service/internal/model"
	"github.com/CoffeeSi/social-network-microservices/content-service/internal/usecase/dto"
)

type PostUseCase interface {
	CreatePost(ctx context.Context, input dto.CreatePostDTO) (model.Post, error)
	GetPostByID(ctx context.Context, id string) (model.Post, error)
	GetFeed(ctx context.Context, input dto.GetFeedDTO) ([]model.Post, int32, error)
	GetUserPosts(ctx context.Context, aID string, input dto.GetFeedDTO) ([]model.Post, int32, error)
	DeletePost(ctx context.Context, id, requestorID string) error
	UpdatePost(ctx context.Context, id, requestorID, content string, mediaURLs []string) (model.Post, error)
}
