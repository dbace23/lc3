package likerepo

import (
	"context"

	"instagram/model"
	"instagram/util/database"
)

type Repo interface {
	Create(ctx context.Context, userID, postID int64) (*model.Like, error)
	ByID(ctx context.Context, id int64) (*model.Like, error)
	DeleteByIDOwner(ctx context.Context, id, ownerID int64) (bool, error)
	ListByPost(ctx context.Context, postID int64) ([]model.Like, error)
	CountByPost(ctx context.Context, postID int64) (int64, error)
}

type repo struct{ db *database.DB }

func New(db *database.DB) Repo { return &repo{db} }

func (r *repo) Create(ctx context.Context, userID, postID int64) (*model.Like, error) {
	var lk model.Like
	err := r.db.Pool.QueryRow(ctx, `
		INSERT INTO likes(user_id, post_id)
		VALUES ($1,$2) RETURNING id, user_id, post_id, created_at`,
		userID, postID,
	).Scan(&lk.ID, &lk.UserID, &lk.PostID, &lk.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &lk, nil
}

func (r *repo) ByID(ctx context.Context, id int64) (*model.Like, error) {
	var lk model.Like
	if err := r.db.Pool.QueryRow(ctx, `
		SELECT 
			id, user_id, post_id, created_at
		FROM 
			likes WHERE id=$1`, id,
	).Scan(&lk.ID, &lk.UserID, &lk.PostID, &lk.CreatedAt); err != nil {
		return nil, err
	}
	return &lk, nil
}

func (r *repo) DeleteByIDOwner(ctx context.Context, id, ownerID int64) (bool, error) {
	cmd, err := r.db.Pool.Exec(ctx, `
		DELETE FROM likes 
		WHERE 
			id=$1 AND user_id=$2`, id, ownerID)
	return cmd.RowsAffected() > 0, err
}

func (r *repo) ListByPost(ctx context.Context, postID int64) ([]model.Like, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT 
			id, user_id, post_id, created_at
		FROM 
			likes WHERE post_id=$1 ORDER BY id DESC`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.Like
	for rows.Next() {
		var lk model.Like
		if err := rows.Scan(&lk.ID, &lk.UserID, &lk.PostID, &lk.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, lk)
	}
	return out, rows.Err()
}

func (r *repo) CountByPost(ctx context.Context, postID int64) (int64, error) {
	var n int64
	err := r.db.Pool.QueryRow(ctx, `SELECT 
										COUNT(distinct ID) 
									FROM 
										likes 
									WHERE 
										post_id=$1`, postID).Scan(&n)
	return n, err
}
