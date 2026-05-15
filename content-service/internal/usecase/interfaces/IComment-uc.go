package interfaces

import (
	"context"

	"github.com/CoffeeSi/social-network-microservices/content-service/internal/model"
	"github.com/CoffeeSi/social-network-microservices/content-service/internal/usecase/dto"
)

type CommentUseCase interface {
	AddComment(ctx context.Context, input dto.CreateCommentDTO) (model.Comment, error)
	ListComments(ctx context.Context, input dto.ListCommentsDTO) ([]model.Comment, int32, error)
	UpdateComment(ctx context.Context, commentID, requestorID, newText string) (model.Comment, error)
	DeleteComment(ctx context.Context, commentID, requestorID string) error
}
