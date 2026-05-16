package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/CoffeeSi/social-network-microservices/content-service/internal/model"
	"github.com/CoffeeSi/social-network-microservices/content-service/internal/model/interfaces"
	"github.com/CoffeeSi/social-network-microservices/content-service/internal/model/interfaces/repos"
)

type LikeUseCase struct {
	likeRepo   repos.LikeRepository
	postRepo   repos.PostRepository
	userClient interfaces.UserServiceClient
	commandBus LikeCommandBus
}

type LikeCommandBus interface {
	RequestToggleLike(ctx context.Context, postID, userID string) (int32, bool, error)
}

func NewLikeUseCase(
	lr repos.LikeRepository,
	pr repos.PostRepository,
	uc interfaces.UserServiceClient) *LikeUseCase {
	return &LikeUseCase{
		likeRepo:   lr,
		postRepo:   pr,
		userClient: uc,
	}
}

func (uc *LikeUseCase) SetCommandBus(bus LikeCommandBus) {
	uc.commandBus = bus
}

func (uc *LikeUseCase) ToggleLike(ctx context.Context, postID, userID string) (int32, bool, error) {
	postID = strings.TrimSpace(postID)
	userID = strings.TrimSpace(userID)
	if postID == "" || userID == "" {
		return 0, false, model.ErrInvalidID
	}
	if uc.commandBus != nil {
		return uc.commandBus.RequestToggleLike(ctx, postID, userID)
	}
	return uc.ProcessToggleLike(ctx, postID, userID)
}

func (uc *LikeUseCase) ProcessToggleLike(ctx context.Context, postID, userID string) (int32, bool, error) {
	postID = strings.TrimSpace(postID)
	userID = strings.TrimSpace(userID)
	if postID == "" || userID == "" {
		return 0, false, model.ErrInvalidID
	}

	exists, err := uc.userClient.UserExists(ctx, userID)
	if err != nil {
		return 0, false, err
	}
	if !exists {
		return 0, false, model.ErrInvalidID
	}

	post, err := uc.postRepo.GetPost(ctx, postID)
	if err != nil {
		return 0, false, err
	}

	like := model.Like{
		PostID:    postID,
		UserID:    userID,
		CreatedAt: time.Now().UTC(),
	}

	if err := uc.likeRepo.LikePost(ctx, like); err == nil {
		if err := uc.postRepo.IncrementLikeCount(ctx, postID, 1); err != nil {
			return 0, false, err
		}
		return post.LikeCount + 1, true, nil
	} else if !errors.Is(err, model.ErrAlreadyLiked) {
		return 0, false, err
	}

	if err := uc.likeRepo.UnlikePost(ctx, postID, userID); err != nil {
		return 0, false, err
	}
	if err := uc.postRepo.IncrementLikeCount(ctx, postID, -1); err != nil {
		return 0, false, err
	}

	if post.LikeCount <= 0 {
		return 0, false, nil
	}
	return post.LikeCount - 1, false, nil
}
