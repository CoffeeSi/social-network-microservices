package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/CoffeeSi/social-network-microservices/content-service/internal/model"
	"github.com/CoffeeSi/social-network-microservices/content-service/internal/repository/dao"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type PostRepo struct {
	col *mongo.Collection
}

func NewPostRepo(col *mongo.Collection) *PostRepo {
	col.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "author_id", Value: 1},
			{Key: "created_at", Value: -1}},
	})
	col.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{{Key: "created_at", Value: -1}},
	})
	return &PostRepo{col: col}
}
func (pr *PostRepo) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	session, err := pr.col.Database().Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx context.Context) (interface{}, error) {
		return nil, fn(sessCtx)
	})
	return err
}
func (r *PostRepo) CreatePost(ctx context.Context, post model.Post) (model.Post, error) {
	postDao, err := dao.FromPostToDao(post)
	if err != nil {
		return model.Post{}, err
	}
	result, err := r.col.InsertOne(ctx, postDao)
	if err != nil {
		return model.Post{}, fmt.Errorf("failed to insert post: %w", err)
	}

	post.ID = result.InsertedID.(bson.ObjectID).Hex()
	return post, nil
}
func (pr *PostRepo) GetPost(ctx context.Context, id string) (model.Post, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return model.Post{}, model.ErrInvalidID
	}

	var p dao.PostDAO
	err = pr.col.FindOne(ctx, bson.M{"_id": objID}).Decode(&p)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Post{}, model.ErrPostNotFound
		}
		return model.Post{}, fmt.Errorf("failed to find post: %w", err)
	}

	return dao.FromDaoToPost(p), nil
}
func (pr *PostRepo) GetUserPosts(ctx context.Context, authorID string, pageSize, page int32) ([]model.Post, int32, error) {
	objAuthorID, err := bson.ObjectIDFromHex(authorID)
	if err != nil {
		return nil, 0, model.ErrInvalidID
	}

	skip := int64((page - 1) * pageSize)

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"author_id": objAuthorID}}},
		{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},
		{{Key: "$facet", Value: bson.D{
			{Key: "data", Value: bson.A{
				bson.D{{Key: "$skip", Value: skip}},
				bson.D{{Key: "$limit", Value: int64(pageSize)}},
			}},
			{Key: "total", Value: bson.A{
				bson.D{{Key: "$count", Value: "count"}},
			}},
		}}},
	}

	return pr.executeAggregationPipeline(ctx, pipeline, pageSize, page)
}
func (pr *PostRepo) GetFeed(ctx context.Context, pageSize, page int32) ([]model.Post, int32, error) {

	skip := int64((page - 1) * pageSize)

	pipeline := mongo.Pipeline{
		{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},
		{{Key: "$facet", Value: bson.D{
			{Key: "data", Value: bson.A{
				bson.D{{Key: "$skip", Value: skip}},
				bson.D{{Key: "$limit", Value: int64(pageSize)}},
			}},
			{Key: "total", Value: bson.A{
				bson.D{{Key: "$count", Value: "count"}},
			}},
		}}},
	}
	return pr.executeAggregationPipeline(ctx, pipeline, pageSize, page)
}
func (pr *PostRepo) DeletePost(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return model.ErrInvalidID
	}

	result, err := pr.col.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}
	if result.DeletedCount == 0 {
		return model.ErrPostNotFound
	}
	return nil
}

func (pr *PostRepo) executeAggregationPipeline(ctx context.Context, pipeline mongo.Pipeline, pageSize, page int32) ([]model.Post, int32, error) {
	cursor, err := pr.col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to aggregate posts: %w", err)
	}
	defer cursor.Close(ctx)

	var result []struct {
		Data  []dao.PostDAO `bson:"data"`
		Total []struct {
			Count int32 `bson:"count"`
		} `bson:"total"`
	}
	if err := cursor.All(ctx, &result); err != nil {
		return nil, 0, err
	}

	if len(result) == 0 {
		return []model.Post{}, 0, nil
	}

	var total int32
	if len(result[0].Total) > 0 {
		total = result[0].Total[0].Count
	}

	posts := make([]model.Post, 0, len(result[0].Data))
	for _, d := range result[0].Data {
		posts = append(posts, dao.FromDaoToPost(d))
	}

	return posts, total, nil
}

func (pr *PostRepo) IncrementCommentCount(ctx context.Context, postID string, amount int32) error {
	objID, err := bson.ObjectIDFromHex(postID)
	if err != nil {
		return model.ErrInvalidID
	}
	result, err := pr.col.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$inc": bson.M{"comment_count": amount}},
	)
	if err != nil {
		return fmt.Errorf("failled to update post comment count: %w", err)
	}
	if result.ModifiedCount == 0 {
		return model.ErrPostNotFound
	}
	return nil
}
func (pr *PostRepo) IncrementLikeCount(ctx context.Context, postID string, amount int32) error {
	objID, err := bson.ObjectIDFromHex(postID)
	if err != nil {
		return model.ErrInvalidID
	}

	result, err := pr.col.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{
		"$inc": bson.M{"like_count": amount},
	})
	if err != nil {
		return fmt.Errorf("failed to update post like count: %w", err)
	}
	if result.ModifiedCount == 0 {
		return model.ErrPostNotFound
	}

	return nil
}
func (pr *PostRepo) UpdatePost(ctx context.Context, id string, content string, mediaURLs []string) (model.Post, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return model.Post{}, model.ErrInvalidID
	}

	update := bson.M{
		"$set": bson.M{
			"content":    content,
			"media_urls": mediaURLs,
			"updated_at": time.Now(),
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedDao dao.PostDAO
	err = pr.col.FindOneAndUpdate(ctx, bson.M{"_id": objID}, update, opts).Decode(&updatedDao)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Post{}, model.ErrPostNotFound
		}
		return model.Post{}, fmt.Errorf("failed to update post: %w", err)
	}

	return dao.FromDaoToPost(updatedDao), nil
}
