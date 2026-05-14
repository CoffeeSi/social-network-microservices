package repos

import (
	"context"

	"github.com/CoffeeSi/social-network-microservices/content-service/internal/model"
)

type LikeRepository interface {
	LikePost(ctx context.Context, like model.Like) error
	UnlikePost(ctx context.Context, postID, userID string) error
}
