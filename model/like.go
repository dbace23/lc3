package model

import "time"
type Like struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	PostID    int64  `json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateLikeReq struct {
	PostID int64 `json:"post_id"`
}
