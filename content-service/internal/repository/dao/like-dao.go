package dao

import (
	"fmt"
	"time"

	"github.com/CoffeeSi/social-network-microservices/content-service/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type LikeDAO struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	PostID    bson.ObjectID `bson:"post_id"`
	UserID    bson.ObjectID `bson:"user_id"`
	CreatedAt time.Time     `bson:"created_at"`
}

func FromLikeToDao(m model.Like) (LikeDAO, error) {
	var pID, uID bson.ObjectID
	var err error

	pID, err = bson.ObjectIDFromHex(m.PostID)
	if err != nil {
		return LikeDAO{}, fmt.Errorf("invalid Post id: %w", err)
	}

	uID, err = bson.ObjectIDFromHex(m.UserID)
	if err != nil {
		return LikeDAO{}, fmt.Errorf("invalid User id: %w", err)
	}

	return LikeDAO{
		PostID:    pID,
		UserID:    uID,
		CreatedAt: m.CreatedAt,
	}, nil
}

func FromDaoToLike(d LikeDAO) model.Like {
	return model.Like{
		UserID:    d.UserID.Hex(),
		PostID:    d.PostID.Hex(),
		CreatedAt: d.CreatedAt,
	}
}
