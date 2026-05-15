package model

import "errors"

var (
	ErrInvalidID        = errors.New("invalid unique identifier format")
	ErrNoFieldsToUpdate = errors.New("at least one field must be provided for update")

	ErrPostNotFound = errors.New("requested post could not be found")
	ErrEmptyPost    = errors.New("cannot create an empty post; text or media is required")

	ErrCommentNotFound = errors.New("requested comment could not be found")
	ErrEmptyComment    = errors.New("comment text cannot be blank")

	ErrAlreadyLiked = errors.New("user has already liked this post")
	ErrLikeNotFound = errors.New("user has not liked this post")

	ErrPermissionDenied = errors.New("permission denied: you are not the author of this content")
)
