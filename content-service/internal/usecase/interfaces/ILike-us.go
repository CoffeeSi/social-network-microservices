package interfaces

import "context"

type LikeUseCase interface {
	LikePost(ctx context.Context, postID, userID string) error
	UnlikePost(ctx context.Context, postID, userID string) error
}
