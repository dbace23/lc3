package model

type Post struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	AuthorID  int64  `json:"author_id"`
	CreatedAt string `json:"created_at"`
}

// model/post.go

// CreatePostReq is the post creation payload
// swagger:model CreatePostReq
type CreatePostReq struct {
	Title   string  `json:"title"`
	Content *string `json:"content"`
}
