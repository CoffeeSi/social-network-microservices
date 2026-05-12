package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/CoffeeSi/social-network-microservices/user-service/internal/model"
	"github.com/CoffeeSi/social-network-microservices/user-service/internal/repository/dao"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type UserRepo struct {
	col *mongo.Collection
}

func NewUserRepo(col *mongo.Collection) *UserRepo {
	col.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	return &UserRepo{col: col}
}

func (ur *UserRepo) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	daouser, err := dao.FromUserToDao(user)
	if err != nil {
		return model.User{}, err
	}
	result, err := ur.col.InsertOne(ctx, daouser)
	if mongo.IsDuplicateKeyError(err) {
		return model.User{}, ErrDuplicateEmail
	}
	if err != nil {
		return model.User{}, err
	}
	user.ID = result.InsertedID.(bson.ObjectID).Hex()
	user.CreatedAt = daouser.CreatedAt
	return user, nil
}

func (ur *UserRepo) GetUsers(ctx context.Context, pageSize, page int32) ([]model.User, int32, error) {
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	if page <= 0 {
		page = 1
	}

	skip := int64((page - 1) * pageSize)

	pipeline := mongo.Pipeline{
		{{Key: "$sort", Value: bson.D{{Key: "_id", Value: 1}}}},
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

	cursor, err := ur.col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to aggregate users: %w", err)
	}
	defer cursor.Close(ctx)

	var result []struct {
		Data  []dao.UserDAO `bson:"data"`
		Total []struct {
			Count int32 `bson:"count"`
		} `bson:"total"`
	}
	if err := cursor.All(ctx, &result); err != nil {
		return nil, 0, fmt.Errorf("failed to decode result: %w", err)
	}

	if len(result) == 0 {
		return []model.User{}, 0, nil
	}

	var total int32
	if len(result[0].Total) > 0 {
		total = result[0].Total[0].Count
	}

	users := make([]model.User, 0, len(result[0].Data))
	for _, u := range result[0].Data {
		users = append(users, dao.FromDaoToUser(u))
	}

	return users, total, nil
}
func (ur *UserRepo) GetUser(ctx context.Context, id string) (model.User, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return model.User{}, ErrInvalidID
	}
	var u dao.UserDAO
	err = ur.col.FindOne(ctx, bson.M{"_id": objID}).Decode(&u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.User{}, ErrUserNotFound
		}
		return model.User{}, err
	}
	return dao.FromDaoToUser(u), nil
}
func (ur *UserRepo) PatchUser(ctx context.Context, id string, updateData model.User) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidID
	}
	updateFields := bson.M{}

	if updateData.FirstName != "" {
		updateFields["first_name"] = updateData.FirstName
	}
	if updateData.LastName != "" {
		updateFields["last_name"] = updateData.LastName
	}
	if updateData.Email != "" {
		email := strings.TrimSpace(strings.ToLower(updateData.Email))
		if match := model.EmailRegex.MatchString(email); !match {
			return ErrInvalidEmail
		}
		updateFields["email"] = email
	}
	if !updateData.DOB.IsZero() {
		updateFields["dob"] = updateData.DOB
	}

	if len(updateFields) == 0 {
		return ErrNoFieldsToUpdate
	}
	result, err := ur.col.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": updateFields},
	)

	if err != nil {
		return ErrOnUpdate
	}

	if result.MatchedCount == 0 {
		return ErrUserNotFound
	}

	return nil
}
func (ur *UserRepo) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	var u dao.UserDAO
	email = strings.TrimSpace(strings.ToLower(email))
	err := ur.col.FindOne(ctx, bson.M{"email": email}).Decode(&u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.User{}, ErrUserNotFound
		}
		return model.User{}, err
	}
	return dao.FromDaoToUserWithPassword(u), nil
}
