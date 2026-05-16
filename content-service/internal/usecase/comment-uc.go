package usecase

import (
	"context"
	"fmt"
	"strings"
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
	input.PostID = strings.TrimSpace(input.PostID)
	input.UserID = strings.TrimSpace(input.UserID)
	input.Text = strings.TrimSpace(input.Text)

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

	_, err = uc.postRepo.GetPost(ctx, input.PostID)
	if err != nil {
		return model.Comment{}, err
	}

	comment := model.Comment{
		PostID:    input.PostID,
		UserID:    input.UserID,
		Text:      input.Text,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	createdComment, err := uc.commentRepo.CreateComment(ctx, comment)
	if err != nil {
		return model.Comment{}, fmt.Errorf("failed to save comment document: %w", err)
	}

	if err := uc.postRepo.IncrementCommentCount(ctx, input.PostID, 1); err != nil {
		return model.Comment{}, fmt.Errorf("failed to increment post comment counter: %w", err)
	}

	return createdComment, nil
}

func (uc *CommentUseCase) ListComments(ctx context.Context, input dto.ListCommentsDTO) ([]model.Comment, int32, error) {
	input.PostID = strings.TrimSpace(input.PostID)
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
	commentID = strings.TrimSpace(commentID)
	requestorID = strings.TrimSpace(requestorID)
	newText = strings.TrimSpace(newText)

	if newText == "" {
		return model.Comment{}, model.ErrEmptyComment
	}
	if commentID == "" || requestorID == "" {
		return model.Comment{}, model.ErrInvalidID
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
	commentID = strings.TrimSpace(commentID)
	requestorID = strings.TrimSpace(requestorID)
	if commentID == "" || requestorID == "" {
		return model.ErrInvalidID
	}

	comment, err := uc.commentRepo.GetComment(ctx, commentID)
	if err != nil {
		return err
	}
	if comment.UserID != requestorID {
		return model.ErrPermissionDenied
	}

	if err := uc.commentRepo.DeleteComment(ctx, commentID); err != nil {
		return fmt.Errorf("failed to delete comment entry: %w", err)
	}
	if err := uc.postRepo.IncrementCommentCount(ctx, comment.PostID, -1); err != nil {
		return fmt.Errorf("failed to decrement post comment counter: %w", err)
	}
	return nil
}
