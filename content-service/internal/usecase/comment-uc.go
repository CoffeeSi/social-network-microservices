package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/CoffeeSi/social-network-microservices/content-service/internal/model"
	"github.com/CoffeeSi/social-network-microservices/content-service/internal/model/interfaces"
	"github.com/CoffeeSi/social-network-microservices/content-service/internal/model/interfaces/repos"
	"github.com/CoffeeSi/social-network-microservices/content-service/internal/usecase/dto"
)

type CommentUseCase struct {
	commentRepo repos.CommentRepository
	postRepo    repos.PostRepository
	userClient  interfaces.UserServiceClient
}

func NewCommentUseCase(
	cr repos.CommentRepository,
	pr repos.PostRepository,
	ucClient interfaces.UserServiceClient,
) *CommentUseCase {
	return &CommentUseCase{
		commentRepo: cr,
		postRepo:    pr,
		userClient:  ucClient,
	}
}

func (uc *CommentUseCase) AddComment(ctx context.Context, input dto.CreateCommentDTO) (model.Comment, error) {
	if input.Text == "" {
		return model.Comment{}, model.ErrEmptyComment
	}
	if input.PostID == "" || input.UserID == "" {
		return model.Comment{}, model.ErrInvalidID
	}

	exists, err := uc.userClient.UserExists(ctx, input.UserID)
	if err != nil {
		return model.Comment{}, fmt.Errorf("failed to verify user status: %w", err)
	}
	if !exists {
		return model.Comment{}, model.ErrInvalidID
	}

	var createdComment model.Comment

	err = uc.postRepo.WithTransaction(ctx, func(txCtx context.Context) error {
		_, err := uc.postRepo.GetPost(txCtx, input.PostID)
		if err != nil {
			return model.ErrPostNotFound
		}

		comment := model.Comment{
			PostID:    input.PostID,
			UserID:    input.UserID,
			Text:      input.Text,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}

		createdComment, err = uc.commentRepo.CreateComment(txCtx, comment)
		if err != nil {
			return fmt.Errorf("failed to save comment document: %w", err)
		}

		err = uc.postRepo.IncrementCommentCount(txCtx, input.PostID, 1)
		if err != nil {
			return fmt.Errorf("failed to increment post comment counter: %w", err)
		}

		return nil
	})
	if err != nil {
		return model.Comment{}, err
	}

	return createdComment, nil
}

func (uc *CommentUseCase) ListComments(ctx context.Context, input dto.ListCommentsDTO) ([]model.Comment, int32, error) {
	if input.PostID == "" {
		return nil, 0, model.ErrInvalidID
	}
	pageSize := input.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	page := input.Page
	if page <= 0 {
		page = 1
	}

	return uc.commentRepo.ListComments(ctx, input.PostID, pageSize, page)
}

func (uc *CommentUseCase) UpdateComment(ctx context.Context, commentID, requestorID, newText string) (model.Comment, error) {

	if newText == "" {
		return model.Comment{}, model.ErrEmptyComment
	}

	existingComment, err := uc.commentRepo.GetComment(ctx, commentID)
	if err != nil {
		return model.Comment{}, err
	}

	if existingComment.UserID != requestorID {
		return model.Comment{}, model.ErrPermissionDenied
	}

	return uc.commentRepo.UpdateComment(ctx, commentID, newText)
}

func (uc *CommentUseCase) DeleteComment(ctx context.Context, commentID, requestorID string) error {
	if commentID == "" || requestorID == "" {
		return model.ErrInvalidID
	}
	return uc.postRepo.WithTransaction(ctx, func(txCtx context.Context) error {
		comment, err := uc.commentRepo.GetComment(txCtx, commentID)
		if err != nil {
			return err
		}
		if comment.UserID != requestorID {
			return model.ErrPermissionDenied
		}

		err = uc.commentRepo.DeleteComment(txCtx, commentID)
		if err != nil {
			return fmt.Errorf("failed to delete comment entry: %w", err)
		}
		err = uc.postRepo.IncrementCommentCount(txCtx, comment.PostID, -1)
		if err != nil {
			return fmt.Errorf("failed to decrement post comment counter: %w", err)
		}
		return nil
	})
}
