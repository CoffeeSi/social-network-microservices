package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/CoffeeSi/social-network-microservices/content-service/internal/model"
	"github.com/CoffeeSi/social-network-microservices/content-service/internal/model/interfaces"
	"github.com/CoffeeSi/social-network-microservices/content-service/internal/model/interfaces/repos"
	"github.com/CoffeeSi/social-network-microservices/content-service/internal/usecase/dto"
)

type PostUseCase struct {
	postRepo   repos.PostRepository
	userClient interfaces.UserServiceClient
}

func NewPostUseCase(repo repos.PostRepository, userClient interfaces.UserServiceClient) *PostUseCase {
	return &PostUseCase{
		postRepo:   repo,
		userClient: userClient,
	}
}

func (uc *PostUseCase) CreatePost(ctx context.Context, input dto.CreatePostDTO) (model.Post, error) {
	aID := strings.TrimSpace(input.AuthorID)
	content := strings.TrimSpace(input.Content)

	if aID == "" {
		return model.Post{}, model.ErrInvalidID
	}
	if content == "" && len(input.MediaURLs) == 0 {
		return model.Post{}, model.ErrEmptyPost
	}

	exists, err := uc.userClient.UserExists(ctx, aID)
	if err != nil {
		return model.Post{}, errors.New("failed to verify author profile availablity")
	}
	if !exists {
		return model.Post{}, errors.New("cannot post: profile does not exist")
	}

	postModel := model.Post{
		AuthorID:     aID,
		Content:      content,
		MediaURLs:    input.MediaURLs,
		LikeCount:    0,
		CommentCount: 0,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	return uc.postRepo.CreatePost(ctx, postModel)
}

func (uc *PostUseCase) GetPostByID(ctx context.Context, id string) (model.Post, error) {
	if strings.TrimSpace(id) == "" {
		return model.Post{}, model.ErrInvalidID
	}
	return uc.postRepo.GetPost(ctx, id)
}
func (uc *PostUseCase) GetFeed(ctx context.Context, input dto.GetFeedDTO) ([]model.Post, int32, error) {
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
	return uc.postRepo.GetFeed(ctx, pageSize, page)
}
func (uc *PostUseCase) GetUserPosts(ctx context.Context, aID string, input dto.GetFeedDTO) ([]model.Post, int32, error) {
	if strings.TrimSpace(aID) == "" {
		return nil, 0, model.ErrInvalidID
	}

	pageSize := input.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	page := input.Page
	if page <= 0 {
		page = 1
	}
	return uc.postRepo.GetUserPosts(ctx, aID, pageSize, page)
}
func (uc *PostUseCase) DeletePost(ctx context.Context, id, requestorID string) error {
	if strings.TrimSpace(id) == "" || strings.TrimSpace(requestorID) == "" {
		return model.ErrInvalidID
	}
	post, err := uc.postRepo.GetPost(ctx, id)
	if err != nil {
		return err
	}
	if post.AuthorID != requestorID {
		return model.ErrPermissionDenied
	}
	return uc.postRepo.DeletePost(ctx, id)
}

func (uc *PostUseCase) UpdatePost(ctx context.Context, id, requestorID, content string, mediaURLs []string) (model.Post, error) {
	id = strings.TrimSpace(id)
	requestorID = strings.TrimSpace(requestorID)
	content = strings.TrimSpace(content)

	if id == "" || requestorID == "" {
		return model.Post{}, model.ErrInvalidID
	}

	if content == "" && len(mediaURLs) == 0 {
		return model.Post{}, model.ErrEmptyPost
	}

	post, err := uc.postRepo.GetPost(ctx, id)
	if err != nil {
		return model.Post{}, err
	}

	if post.AuthorID != requestorID {
		return model.Post{}, model.ErrPermissionDenied
	}

	return uc.postRepo.UpdatePost(ctx, id, content, mediaURLs)
}
