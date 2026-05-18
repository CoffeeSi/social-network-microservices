package usecase

import (
	"context"
	"testing"

	"github.com/CoffeeSi/social-network-microservices/content-service/internal/model"
	"github.com/CoffeeSi/social-network-microservices/content-service/internal/usecase/dto"
)

func TestAddCommentUsesTransactionAndIncrementsCount(t *testing.T) {
	ctx := context.Background()
	postRepo := newFakePostRepo()
	postRepo.posts["post1"] = model.Post{ID: "post1", AuthorID: "user1"}
	commentRepo := newFakeCommentRepo()
	uc := NewCommentUseCase(commentRepo, postRepo, fakeUserClient{exists: true})

	comment, err := uc.AddComment(ctx, dto.CreateCommentDTO{
		PostID: "post1",
		UserID: "user1",
		Text:   "  hello  ",
	})
	if err != nil {
		t.Fatalf("AddComment returned error: %v", err)
	}

	if postRepo.transactionCount != 1 {
		t.Fatalf("expected one transaction, got %d", postRepo.transactionCount)
	}
	if postRepo.commentIncrements["post1"] != 1 {
		t.Fatalf("expected comment count increment 1, got %d", postRepo.commentIncrements["post1"])
	}
	if comment.ID == "" {
		t.Fatal("expected created comment id")
	}
	if comment.Text != "hello" {
		t.Fatalf("expected trimmed text, got %q", comment.Text)
	}
}

func TestDeleteCommentUsesTransactionAndDecrementsCount(t *testing.T) {
	ctx := context.Background()
	postRepo := newFakePostRepo()
	commentRepo := newFakeCommentRepo()
	commentRepo.comments["comment1"] = model.Comment{
		ID:     "comment1",
		PostID: "post1",
		UserID: "user1",
	}
	uc := NewCommentUseCase(commentRepo, postRepo, fakeUserClient{exists: true})

	if err := uc.DeleteComment(ctx, "comment1", "user1"); err != nil {
		t.Fatalf("DeleteComment returned error: %v", err)
	}

	if postRepo.transactionCount != 1 {
		t.Fatalf("expected one transaction, got %d", postRepo.transactionCount)
	}
	if postRepo.commentIncrements["post1"] != -1 {
		t.Fatalf("expected comment count decrement -1, got %d", postRepo.commentIncrements["post1"])
	}
	if _, ok := commentRepo.comments["comment1"]; ok {
		t.Fatal("expected comment to be deleted")
	}
}

func TestProcessToggleLikeUsesTransactionForLikeAndUnlike(t *testing.T) {
	ctx := context.Background()
	postRepo := newFakePostRepo()
	postRepo.posts["post1"] = model.Post{ID: "post1", AuthorID: "author1"}
	likeRepo := newFakeLikeRepo()
	uc := NewLikeUseCase(likeRepo, postRepo, fakeUserClient{exists: true})

	count, liked, err := uc.ProcessToggleLike(ctx, "post1", "user1")
	if err != nil {
		t.Fatalf("ProcessToggleLike like returned error: %v", err)
	}
	if !liked || count != 1 {
		t.Fatalf("expected liked=true count=1, got liked=%v count=%d", liked, count)
	}
	if postRepo.transactionCount != 1 {
		t.Fatalf("expected one transaction after like, got %d", postRepo.transactionCount)
	}
	if postRepo.likeIncrements["post1"] != 1 {
		t.Fatalf("expected like increment 1, got %d", postRepo.likeIncrements["post1"])
	}

	count, liked, err = uc.ProcessToggleLike(ctx, "post1", "user1")
	if err != nil {
		t.Fatalf("ProcessToggleLike unlike returned error: %v", err)
	}
	if liked || count != 0 {
		t.Fatalf("expected liked=false count=0, got liked=%v count=%d", liked, count)
	}
	if postRepo.transactionCount != 2 {
		t.Fatalf("expected two transactions after unlike, got %d", postRepo.transactionCount)
	}
	if postRepo.likeIncrements["post1"] != 0 {
		t.Fatalf("expected net like increment 0, got %d", postRepo.likeIncrements["post1"])
	}
}

type fakeUserClient struct {
	exists bool
	err    error
}

func (f fakeUserClient) UserExists(ctx context.Context, userID string) (bool, error) {
	return f.exists, f.err
}

type fakePostRepo struct {
	posts             map[string]model.Post
	commentIncrements map[string]int32
	likeIncrements    map[string]int32
	transactionCount  int
}

func newFakePostRepo() *fakePostRepo {
	return &fakePostRepo{
		posts:             map[string]model.Post{},
		commentIncrements: map[string]int32{},
		likeIncrements:    map[string]int32{},
	}
}

func (f *fakePostRepo) CreatePost(ctx context.Context, post model.Post) (model.Post, error) {
	if post.ID == "" {
		post.ID = "post-created"
	}
	f.posts[post.ID] = post
	return post, nil
}

func (f *fakePostRepo) GetPost(ctx context.Context, id string) (model.Post, error) {
	post, ok := f.posts[id]
	if !ok {
		return model.Post{}, model.ErrPostNotFound
	}
	return post, nil
}

func (f *fakePostRepo) DeletePost(ctx context.Context, id string) error {
	delete(f.posts, id)
	return nil
}

func (f *fakePostRepo) UpdatePost(ctx context.Context, id string, content string, mediaURLs []string) (model.Post, error) {
	post, ok := f.posts[id]
	if !ok {
		return model.Post{}, model.ErrPostNotFound
	}
	post.Content = content
	post.MediaURLs = mediaURLs
	f.posts[id] = post
	return post, nil
}

func (f *fakePostRepo) GetFeed(ctx context.Context, pageSize, page int32) ([]model.Post, int32, error) {
	return nil, 0, nil
}

func (f *fakePostRepo) GetUserPosts(ctx context.Context, authorID string, pageSize, page int32) ([]model.Post, int32, error) {
	return nil, 0, nil
}

func (f *fakePostRepo) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	f.transactionCount++
	return fn(ctx)
}

func (f *fakePostRepo) IncrementCommentCount(ctx context.Context, postID string, amount int32) error {
	post, ok := f.posts[postID]
	if !ok && len(f.posts) > 0 {
		return model.ErrPostNotFound
	}
	post.CommentCount += amount
	f.posts[postID] = post
	f.commentIncrements[postID] += amount
	return nil
}

func (f *fakePostRepo) IncrementLikeCount(ctx context.Context, postID string, amount int32) error {
	post, ok := f.posts[postID]
	if !ok {
		return model.ErrPostNotFound
	}
	post.LikeCount += amount
	f.posts[postID] = post
	f.likeIncrements[postID] += amount
	return nil
}

type fakeCommentRepo struct {
	comments map[string]model.Comment
	nextID   int
}

func newFakeCommentRepo() *fakeCommentRepo {
	return &fakeCommentRepo{comments: map[string]model.Comment{}}
}

func (f *fakeCommentRepo) CreateComment(ctx context.Context, comment model.Comment) (model.Comment, error) {
	f.nextID++
	comment.ID = "comment-created"
	f.comments[comment.ID] = comment
	return comment, nil
}

func (f *fakeCommentRepo) GetComment(ctx context.Context, id string) (model.Comment, error) {
	comment, ok := f.comments[id]
	if !ok {
		return model.Comment{}, model.ErrCommentNotFound
	}
	return comment, nil
}

func (f *fakeCommentRepo) ListComments(ctx context.Context, postID string, pageSize, page int32) ([]model.Comment, int32, error) {
	return nil, 0, nil
}

func (f *fakeCommentRepo) UpdateComment(ctx context.Context, id string, text string) (model.Comment, error) {
	comment, ok := f.comments[id]
	if !ok {
		return model.Comment{}, model.ErrCommentNotFound
	}
	comment.Text = text
	f.comments[id] = comment
	return comment, nil
}

func (f *fakeCommentRepo) DeleteComment(ctx context.Context, id string) error {
	if _, ok := f.comments[id]; !ok {
		return model.ErrCommentNotFound
	}
	delete(f.comments, id)
	return nil
}

type fakeLikeRepo struct {
	likes map[string]model.Like
}

func newFakeLikeRepo() *fakeLikeRepo {
	return &fakeLikeRepo{likes: map[string]model.Like{}}
}

func (f *fakeLikeRepo) LikePost(ctx context.Context, like model.Like) error {
	key := like.PostID + ":" + like.UserID
	if _, ok := f.likes[key]; ok {
		return model.ErrAlreadyLiked
	}
	f.likes[key] = like
	return nil
}

func (f *fakeLikeRepo) UnlikePost(ctx context.Context, postID, userID string) error {
	key := postID + ":" + userID
	if _, ok := f.likes[key]; !ok {
		return model.ErrLikeNotFound
	}
	delete(f.likes, key)
	return nil
}

func (f *fakeLikeRepo) IsLiked(ctx context.Context, postID, userID string) (bool, error) {
	_, ok := f.likes[postID+":"+userID]
	return ok, nil
}
