package model

import "time"

type Post struct {
	ID           string
	AuthorID     string
	Content      string
	MediaURLs    []string
	LikeCount    int32
	CommentCount int32
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
