package model

import "time"

type Like struct {
	PostID    string
	UserID    string
	CreatedAt time.Time
}
