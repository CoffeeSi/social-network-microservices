package usecase

import (
	"context"
	"time"

	"github.com/CoffeeSi/social-network-microservices/content-service/internal/model"
	"github.com/CoffeeSi/social-network-microservices/content-service/internal/model/interfaces"
	"github.com/CoffeeSi/social-network-microservices/content-service/internal/model/interfaces/repos"
)

type LikeUseCase struct {
	likeRepo   repos.LikeRepository
	postRepo   repos.PostRepository
	userClient interfaces.UserServiceClient
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

func (uc *LikeUseCase) LikePost(ctx context.Context, postID, userID string) error {
	if postID == "" || userID == "" {
		return model.ErrInvalidID
	}

	exists, err := uc.userClient.UserExists(ctx, userID)
	if err != nil || !exists {
		return model.ErrInvalidID
	}

	return uc.postRepo.WithTransaction(ctx, func(txCtx context.Context) error {
		_, err := uc.postRepo.GetPost(txCtx, postID)
		if err != nil {
			return model.ErrPostNotFound
		}

		like := model.Like{
			PostID:    postID,
			UserID:    userID,
			CreatedAt: time.Now().UTC(),
		}
		err = uc.likeRepo.LikePost(txCtx, like)
		if err != nil {
			return err
		}
		return uc.postRepo.IncrementLikeCount(txCtx, postID, 1)
	})
}

func (uc *LikeUseCase) UnlikePost(ctx context.Context, postID, userID string) error {
	if postID == "" || userID == "" {
		return model.ErrInvalidID
	}

	return uc.postRepo.WithTransaction(ctx, func(txCtx context.Context) error {
		err := uc.likeRepo.UnlikePost(txCtx, postID, userID)
		if err != nil {
			return err
		}
		return uc.postRepo.IncrementLikeCount(txCtx, postID, -1)
	})
}
