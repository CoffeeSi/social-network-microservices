package dto

type CreateCommentDTO struct {
	PostID string `json:"post_id"`
	UserID string `json:"user_id"`
	Text   string `json:"text"`
}
type ListCommentsDTO struct {
	PostID   string `json:"post_id"`
	Page     int32  `json:"page"`
	PageSize int32  `json:"page_size"`
}
