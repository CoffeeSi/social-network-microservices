package dto

type CreatePostDTO struct {
	AuthorID  string   `json:"author_id"`
	Content   string   `json:"content"`
	MediaURLs []string `json:"media_urls"`
}

type GetFeedDTO struct {
	Page     int32 `json:"page"`
	PageSize int32 `json:"page_size"`
}
