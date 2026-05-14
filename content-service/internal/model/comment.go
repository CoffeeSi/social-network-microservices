package model

import "time"

type Comment struct {
	ID        string
	PostID    string
	UserID    string
	Text      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
