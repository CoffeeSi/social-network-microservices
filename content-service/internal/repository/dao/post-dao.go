package dao

import (
	"fmt"
	"time"

	"github.com/CoffeeSi/social-network-microservices/content-service/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type PostDAO struct {
	ID           bson.ObjectID `bson:"_id,omitempty"`
	AuthorID     bson.ObjectID `bson:"author_id"`
	Content      string        `bson:"content"`
	MediaURLs    []string      `bson:"media_urls"`
	LikeCount    int32         `bson:"like_count"`
	CommentCount int32         `bson:"comment_count"`
	CreatedAt    time.Time     `bson:"created_at"`
	UpdatedAt    time.Time     `bson:"updated_at"`
}

func FromPostToDao(m model.Post) (PostDAO, error) {
	var id, authID bson.ObjectID
	var err error

	if m.ID != "" {
		id, err = bson.ObjectIDFromHex(m.ID)
		if err != nil {
			return PostDAO{}, fmt.Errorf("invalid object id: %w", err)
		}
	}
	authID, err = bson.ObjectIDFromHex(m.AuthorID)
	if err != nil {
		return PostDAO{}, err
	}

	return PostDAO{
		ID:           id,
		AuthorID:     authID,
		Content:      m.Content,
		MediaURLs:    m.MediaURLs,
		LikeCount:    m.LikeCount,
		CommentCount: m.CommentCount,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}, nil
}

func FromDaoToPost(d PostDAO) model.Post {
	return model.Post{
		ID:           d.ID.Hex(),
		AuthorID:     d.AuthorID.Hex(),
		Content:      d.Content,
		MediaURLs:    d.MediaURLs,
		LikeCount:    d.LikeCount,
		CommentCount: d.CommentCount,
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
	}
}
