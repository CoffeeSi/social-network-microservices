package events

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/CoffeeSi/social-network-microservices/content-service/internal/model"
	"github.com/nats-io/nats.go"
)

const likeToggleSubject = "content.likes.toggle"

type LikeCommandHandler interface {
	ProcessToggleLike(ctx context.Context, postID, userID string) (int32, bool, error)
}

type LikeQueue struct {
	nc      *nats.Conn
	handler LikeCommandHandler
	sub     *nats.Subscription
	timeout time.Duration
}

type likeCommand struct {
	PostID string `json:"post_id"`
	UserID string `json:"user_id"`
}

type likeCommandResponse struct {
	NewLikeCount int32  `json:"new_like_count"`
	IsLiked      bool   `json:"is_liked"`
	Error        string `json:"error,omitempty"`
}

func NewLikeQueue(nc *nats.Conn, handler LikeCommandHandler) *LikeQueue {
	return &LikeQueue{
		nc:      nc,
		handler: handler,
		timeout: 5 * time.Second,
	}
}

func (q *LikeQueue) Start() error {
	sub, err := q.nc.QueueSubscribe(likeToggleSubject, "content-like-workers", q.handleMessage)
	if err != nil {
		return err
	}
	q.sub = sub
	return q.nc.Flush()
}

func (q *LikeQueue) Stop() {
	if q.sub != nil {
		_ = q.sub.Drain()
	}
}

func (q *LikeQueue) RequestToggleLike(ctx context.Context, postID, userID string) (int32, bool, error) {
	payload, err := json.Marshal(likeCommand{PostID: postID, UserID: userID})
	if err != nil {
		return 0, false, err
	}

	reqCtx, cancel := context.WithTimeout(ctx, q.timeout)
	defer cancel()

	msg, err := q.nc.RequestWithContext(reqCtx, likeToggleSubject, payload)
	if err != nil {
		return 0, false, err
	}

	var res likeCommandResponse
	if err := json.Unmarshal(msg.Data, &res); err != nil {
		return 0, false, err
	}
	if res.Error != "" {
		return 0, false, errorFromCode(res.Error)
	}

	return res.NewLikeCount, res.IsLiked, nil
}

func (q *LikeQueue) handleMessage(msg *nats.Msg) {
	var cmd likeCommand
	if err := json.Unmarshal(msg.Data, &cmd); err != nil {
		q.respond(msg, likeCommandResponse{Error: "invalid_id"})
		return
	}

	count, liked, err := q.handler.ProcessToggleLike(context.Background(), cmd.PostID, cmd.UserID)
	if err != nil {
		q.respond(msg, likeCommandResponse{Error: codeFromError(err)})
		return
	}

	eventSubject := "posts.unliked"
	if liked {
		eventSubject = "posts.liked"
	}
	_ = q.publishLikeEvent(eventSubject, cmd.PostID, cmd.UserID, count)
	q.respond(msg, likeCommandResponse{NewLikeCount: count, IsLiked: liked})
}

func (q *LikeQueue) respond(msg *nats.Msg, res likeCommandResponse) {
	data, _ := json.Marshal(res)
	_ = msg.Respond(data)
}

func (q *LikeQueue) publishLikeEvent(subject, postID, userID string, count int32) error {
	payload := struct {
		EventType    string `json:"event_type"`
		OccurredAt   string `json:"occurred_at"`
		PostID       string `json:"post_id"`
		UserID       string `json:"user_id"`
		NewLikeCount int32  `json:"new_like_count"`
	}{
		EventType:    subject,
		OccurredAt:   time.Now().UTC().Format(time.RFC3339),
		PostID:       postID,
		UserID:       userID,
		NewLikeCount: count,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return q.nc.Publish(subject, data)
}

func codeFromError(err error) string {
	switch {
	case errors.Is(err, model.ErrInvalidID):
		return "invalid_id"
	case errors.Is(err, model.ErrPostNotFound):
		return "post_not_found"
	case errors.Is(err, model.ErrLikeNotFound):
		return "like_not_found"
	default:
		return "internal"
	}
}

func errorFromCode(code string) error {
	switch code {
	case "invalid_id":
		return model.ErrInvalidID
	case "post_not_found":
		return model.ErrPostNotFound
	case "like_not_found":
		return model.ErrLikeNotFound
	default:
		return errors.New("like command processing failed")
	}
}
