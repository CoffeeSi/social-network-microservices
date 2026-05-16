package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func RunMigrations(ctx context.Context, database *mongo.Database) error {
	if err := createPostIndexes(ctx, database.Collection("posts")); err != nil {
		return fmt.Errorf("create posts indexes: %w", err)
	}
	if err := createCommentIndexes(ctx, database.Collection("comments")); err != nil {
		return fmt.Errorf("create comments indexes: %w", err)
	}
	if err := createLikeIndexes(ctx, database.Collection("likes")); err != nil {
		return fmt.Errorf("create likes indexes: %w", err)
	}
	return nil
}

func createPostIndexes(ctx context.Context, col *mongo.Collection) error {
	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "author_id", Value: 1},
				{Key: "created_at", Value: -1},
			},
			Options: options.Index().SetName("idx_posts_author_created_at"),
		},
		{
			Keys:    bson.D{{Key: "created_at", Value: -1}},
			Options: options.Index().SetName("idx_posts_created_at"),
		},
	})
	return err
}

func createCommentIndexes(ctx context.Context, col *mongo.Collection) error {
	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "post_id", Value: 1},
				{Key: "created_at", Value: 1},
			},
			Options: options.Index().SetName("idx_comments_post_created_at"),
		},
	})
	return err
}

func createLikeIndexes(ctx context.Context, col *mongo.Collection) error {
	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "post_id", Value: 1},
				{Key: "user_id", Value: 1},
			},
			Options: options.Index().SetName("idx_likes_post_user_unique").SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "post_id", Value: 1}},
			Options: options.Index().SetName("idx_likes_post"),
		},
		{
			Keys:    bson.D{{Key: "user_id", Value: 1}},
			Options: options.Index().SetName("idx_likes_user"),
		},
	})
	return err
}
