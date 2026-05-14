package dao

import (
	"fmt"
	"time"

	"github.com/CoffeeSi/social-network-microservices/content-service/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type CommentDAO struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	PostID    bson.ObjectID `bson:"post_id"`
	UserID    bson.ObjectID `bson:"user_id"`
	Text      string        `bson:"text"`
	CreatedAt time.Time     `bson:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at"`
}

func FromCommentToDao(m model.Comment) (CommentDAO, error) {
	var id, pID, uID bson.ObjectID
	var err error

	if m.ID != "" {
		id, err = bson.ObjectIDFromHex(m.ID)
		if err != nil {
			return CommentDAO{}, fmt.Errorf("invalid Comment id: %w", err)
		}
	}

	pID, err = bson.ObjectIDFromHex(m.PostID)
	if err != nil {
		return CommentDAO{}, fmt.Errorf("invalid Post id: %w", err)
	}

	uID, err = bson.ObjectIDFromHex(m.UserID)
	if err != nil {
		return CommentDAO{}, fmt.Errorf("invalid User id: %w", err)
	}

	return CommentDAO{
		ID:        id,
		PostID:    pID,
		UserID:    uID,
		Text:      m.Text,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}, nil

}

func FromDaoToComment(d CommentDAO) model.Comment {
	return model.Comment{
		ID:        d.ID.Hex(),
		PostID:    d.PostID.Hex(),
		UserID:    d.UserID.Hex(),
		Text:      d.Text,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}
