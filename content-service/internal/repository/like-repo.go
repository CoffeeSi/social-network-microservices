package repository

import (
	"context"
	"fmt"

	"github.com/CoffeeSi/social-network-microservices/content-service/internal/model"
	"github.com/CoffeeSi/social-network-microservices/content-service/internal/repository/dao"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type LikeRepo struct {
	col *mongo.Collection
}

func NewLikeRepo(col *mongo.Collection) *LikeRepo {
	return &LikeRepo{col: col}
}

func (r *LikeRepo) LikePost(ctx context.Context, like model.Like) error {
	likeDao, err := dao.FromLikeToDao(like)
	if err != nil {
		return err
	}
	_, err = r.col.InsertOne(ctx, likeDao)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return model.ErrAlreadyLiked
		}
		return fmt.Errorf("failed to insert like record: %w", err)
	}
	return nil
}

func (r *LikeRepo) UnlikePost(ctx context.Context, postID, userID string) error {
	pID, err := bson.ObjectIDFromHex(postID)
	if err != nil {
		return model.ErrInvalidID
	}
	uID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return model.ErrInvalidID
	}
	result, err := r.col.DeleteOne(ctx, bson.M{
		"post_id": pID,
		"user_id": uID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete liek record: %w", err)
	}
	if result.DeletedCount == 0 {
		return model.ErrLikeNotFound
	}
	return nil
}
