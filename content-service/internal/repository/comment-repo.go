package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/CoffeeSi/social-network-microservices/content-service/internal/model"
	"github.com/CoffeeSi/social-network-microservices/content-service/internal/repository/dao"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type CommentRepo struct {
	col *mongo.Collection
}

func NewCommentRepo(col *mongo.Collection) *CommentRepo {
	col.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "post_id", Value: 1},
			{Key: "created_at", Value: 1},
		},
	})
	return &CommentRepo{col: col}
}

func (r *CommentRepo) CreateComment(ctx context.Context, comment model.Comment) (model.Comment, error) {
	commentDao, err := dao.FromCommentToDao(comment)
	if err != nil {
		return model.Comment{}, err
	}

	result, err := r.col.InsertOne(ctx, commentDao)
	if err != nil {
		return model.Comment{}, fmt.Errorf("failed to insert comment: %w", err)
	}

	comment.ID = result.InsertedID.(bson.ObjectID).Hex()
	return comment, nil
}
func (r *CommentRepo) GetComment(ctx context.Context, id string) (model.Comment, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return model.Comment{}, model.ErrInvalidID
	}
	var c dao.CommentDAO
	err = r.col.FindOne(ctx, bson.M{"_id": objID}).Decode(&c)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Comment{}, model.ErrCommentNotFound
		}
		return model.Comment{}, fmt.Errorf("failed to find comment: %w", err)
	}
	return dao.FromDaoToComment(c), nil
}

func (r *CommentRepo) ListComments(ctx context.Context, pID string, pageSize, page int32) ([]model.Comment, int32, error) {
	objId, err := bson.ObjectIDFromHex(pID)
	if err != nil {
		return nil, 0, model.ErrInvalidID
	}
	skip := int64((page - 1) * pageSize)

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"post_id": objId}}},
		{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: 1}}}},
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
	cursor, err := r.col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list comments: %w", err)
	}
	defer cursor.Close(ctx)

	var result []struct {
		Data  []dao.CommentDAO `bson:"data"`
		Total []struct {
			Count int32 `bson:"count"`
		} `bson:"total"`
	}

	if err := cursor.All(ctx, &result); err != nil {
		return nil, 0, err
	}

	if len(result) == 0 || len(result[0].Data) == 0 {
		return []model.Comment{}, 0, nil
	}

	var total int32
	if len(result[0].Total) > 0 {
		total = result[0].Total[0].Count
	}

	comments := make([]model.Comment, 0, len(result[0].Data))
	for _, d := range result[0].Data {
		comments = append(comments, dao.FromDaoToComment(d))
	}

	return comments, total, nil
}

func (r *CommentRepo) DeleteComment(ctx context.Context, id string) error {
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return model.ErrInvalidID
	}
	result, err := r.col.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return fmt.Errorf("faile to delete comment: %w", err)
	}
	if result.DeletedCount == 0 {
		return model.ErrCommentNotFound
	}
	return nil
}

func (r *CommentRepo) UpdateComment(ctx context.Context, id string, text string) (model.Comment, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return model.Comment{}, model.ErrInvalidID
	}

	update := bson.M{
		"$set": bson.M{
			"text":       text,
			"updated_at": bson.M{"$literal": "NOW"},
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedDao dao.CommentDAO
	err = r.col.FindOneAndUpdate(ctx, bson.M{"_id": objID}, update, opts).Decode(&updatedDao)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Comment{}, model.ErrCommentNotFound
		}
		return model.Comment{}, fmt.Errorf("failed to update comment: %w", err)
	}

	return dao.FromDaoToComment(updatedDao), nil
}
