package usecase

import (
	"context"
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
		count, liked, err := uc.commandBus.RequestToggleLike(ctx, postID, userID)
		if err == nil {
			return count, liked, nil
		}
		// Fallback when NATS request/reply is unavailable (e.g. timeout during local dev).
		return uc.ProcessToggleLike(ctx, postID, userID)
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

	alreadyLiked, err := uc.likeRepo.IsLiked(ctx, postID, userID)
	if err != nil {
		return 0, false, err
	}

	var count int32
	var liked bool
	err = uc.postRepo.WithTransaction(ctx, func(txCtx context.Context) error {
		if !alreadyLiked {
			like := model.Like{
				PostID:    postID,
				UserID:    userID,
				CreatedAt: time.Now().UTC(),
			}

			if err := uc.likeRepo.LikePost(txCtx, like); err != nil {
				return err
			}
			if err := uc.postRepo.IncrementLikeCount(txCtx, postID, 1); err != nil {
				return err
			}

			count = post.LikeCount + 1
			liked = true
			return nil
		}

		if err := uc.likeRepo.UnlikePost(txCtx, postID, userID); err != nil {
			return err
		}
		if err := uc.postRepo.IncrementLikeCount(txCtx, postID, -1); err != nil {
			return err
		}
		if post.LikeCount <= 0 {
			count = 0
		} else {
			count = post.LikeCount - 1
		}
		liked = false
		return nil
	})
	if err != nil {
		return 0, false, err
	}

	return count, liked, nil
}
