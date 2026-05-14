package repos

import (
	"context"

	"github.com/CoffeeSi/social-network-microservices/content-service/internal/model"
)

type PostEventPublisher interface {
	PublishPostCreated(ctx context.Context, post model.Post) error
}
