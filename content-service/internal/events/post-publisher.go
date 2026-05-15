package events

import (
	"context"
	"encoding/json"
	"time"

	"github.com/CoffeeSi/social-network-microservices/content-service/internal/model"
	"github.com/CoffeeSi/social-network-microservices/content-service/internal/model/interfaces/repos"
	"github.com/nats-io/nats.go"
)

type postPublisher struct {
	nc *nats.Conn
}

func NewPostPublisher(nc *nats.Conn) repos.PostEventPublisher {
	return &postPublisher{nc: nc}
}
func (p *postPublisher) PublishPostCreated(ctx context.Context, post model.Post) error {
	payload := struct {
		EventType  string `json:"event_type"`
		OccurredAt string `json:"occurred_at"`
		PostID     string `json:"post_id"`
		AuthorID   string `json:"author_id"`
		Content    string `json:"content"`
		HasMedia   bool   `json:"has_media"`
	}{
		EventType:  "posts.created",
		OccurredAt: time.Now().Format(time.RFC3339),
		PostID:     post.ID,
		AuthorID:   post.AuthorID,
		Content:    post.Content,
		HasMedia:   len(post.MediaURLs) > 0,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return p.nc.Publish("posts.created", data)
}
