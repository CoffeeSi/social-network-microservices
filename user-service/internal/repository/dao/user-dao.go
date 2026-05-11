package dao

import (
	"fmt"
	"time"

	"github.com/CoffeeSi/social-network-microservices/user-service/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserDAO struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	FirstName string             `bson:"first_name"`
	LastName  string             `bson:"last_name"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
	DOB       time.Time          `bson:"dob"`
	IsActive  bool               `bson:"is_active"`
}

func FromUserToDao(u model.User) (UserDAO, error) {
	var objId primitive.ObjectID
	if u.ID != "" {
		var err error
		objId, err = primitive.ObjectIDFromHex(u.ID)
		if err != nil {
			return UserDAO{}, fmt.Errorf("invalid object id: %w", err)
		}
	}
	layout := "02-01-2006"
	t, err := time.Parse(layout, u.DOB)
	if err != nil {
		return UserDAO{}, fmt.Errorf("invalid date format: %w", err)
	}

	return UserDAO{
		ID:        objId,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Password:  u.Password,
		DOB:       t,
		IsActive:  u.IsActive,
	}, nil
}

func FromDaoToUser(u UserDAO) model.User {
	return model.User{
		ID:        u.ID.Hex(),
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Password:  "",
		DOB:       u.DOB.Format("02-01-2006"),
		IsActive:  u.IsActive,
	}
}

func FromDaoToUserWithPassword(u UserDAO) model.User {
	user := FromDaoToUser(u)
	user.Password = u.Password
	return user
}
