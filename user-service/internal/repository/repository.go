package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/CoffeeSi/social-network-microservices/user-service/internal/model"
	"github.com/CoffeeSi/social-network-microservices/user-service/internal/repository/dao"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	user.ID = result.InsertedID.(primitive.ObjectID).Hex()

	return user, nil
}

func (ur *UserRepo) GetUsers(ctx context.Context) ([]model.User, error) {
	cursor, err := ur.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find users: %w", err)
	}
	defer cursor.Close(ctx)

	var daos []dao.UserDAO
	if err := cursor.All(ctx, &daos); err != nil {
		return nil, fmt.Errorf("failed to decode users: %w", err)
	}

	result := make([]model.User, 0, len(daos))
	for _, u := range daos {
		userModel := dao.FromDaoToUser(u)
		result = append(result, userModel)
	}

	return result, nil
}
func (ur *UserRepo) GetUser(ctx context.Context, id string) (model.User, error) {
	objID, err := primitive.ObjectIDFromHex(id)
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
	objID, err := primitive.ObjectIDFromHex(id)
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
	if updateData.DOB != "" {
		t, err := time.Parse(model.DateLayout, updateData.DOB)
		if err != nil {
			return fmt.Errorf("invalid date format: %w", err)
		}
		updateFields["dob"] = t
	}

	if len(updateFields) == 0 {
		return fmt.Errorf("no fields to update")
	}
	result, err := ur.col.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": updateFields},
	)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return ErrUserNotFound
	}

	return nil
}
func (ur *UserRepo) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	var u dao.UserDAO
	err := ur.col.FindOne(ctx, bson.M{"email": email}).Decode(&u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.User{}, ErrUserNotFound
		}
		return model.User{}, err
	}
	return dao.FromDaoToUserWithPassword(u), nil
}
