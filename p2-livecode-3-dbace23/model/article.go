package model

type Article struct {
	ID         int64  `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	AuthorID   int64  `json:"author_id"`
	CategoryID int64  `json:"category_id"`
	CreatedAt  string `json:"created_at"`
}

type CreateArticleReq struct {
	Title      string  `json:"title"`
	Content    *string `json:"content"`
	CategoryID int64   `json:"category_id"`
}

type ArticleDetail struct {
	Article
	LikesCount int64   `json:"likes_count,omitempty"`
	LikerIDs   []int64 `json:"liker_ids,omitempty"`
}
