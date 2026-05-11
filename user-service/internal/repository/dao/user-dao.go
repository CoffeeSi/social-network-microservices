package dao

import (
	"fmt"
	"time"

	"github.com/CoffeeSi/social-network-microservices/user-service/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserDAO struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	FirstName string        `bson:"first_name"`
	LastName  string        `bson:"last_name"`
	Email     string        `bson:"email"`
	Password  string        `bson:"password"`
	DOB       time.Time     `bson:"dob"`
	IsActive  bool          `bson:"is_active"`
	CreatedAt time.Time     `bson:"created_at"`
}

func FromUserToDao(u model.User) (UserDAO, error) {
	var objId bson.ObjectID
	if u.ID != "" {
		var err error
		objId, err = bson.ObjectIDFromHex(u.ID)
		if err != nil {
			return UserDAO{}, fmt.Errorf("invalid object id: %w", err)
		}
	}

	return UserDAO{
		ID:        objId,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Password:  u.Password,
		DOB:       u.DOB,
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
		DOB:       u.DOB,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
	}
}

func FromDaoToUserWithPassword(u UserDAO) model.User {
	user := FromDaoToUser(u)
	user.Password = u.Password
	return user
}
