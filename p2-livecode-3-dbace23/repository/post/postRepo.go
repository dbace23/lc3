package postrepo

import (
	"context"

	"instagram/model"
	"instagram/util/database"
)

type Repo interface {
	Create(ctx context.Context, p *model.Post) error
	All(ctx context.Context) ([]model.Post, error)
	ByID(ctx context.Context, id int64) (*model.Post, error)
	DeleteByIDOwner(ctx context.Context, id, ownerID int64) (bool, error)
}

type repo struct{ db *database.DB }

func New(db *database.DB) Repo { return &repo{db} }

func (r *repo) Create(ctx context.Context, p *model.Post) error {
	return r.db.Pool.QueryRow(ctx, `
		INSERT INTO posts(title, content, author_id)
		VALUES ($1,$2,$3) RETURNING id, created_at`,
		p.Title, p.Content, p.AuthorID,
	).Scan(&p.ID, &p.CreatedAt)
}

func (r *repo) All(ctx context.Context) ([]model.Post, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT 
			id, title, content, author_id, created_at
		FROM 
			posts ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.Post
	for rows.Next() {
		var p model.Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.AuthorID, &p.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (r *repo) ByID(ctx context.Context, id int64) (*model.Post, error) {
	var p model.Post
	if err := r.db.Pool.QueryRow(ctx, `
		SELECT 
			id, title, content, author_id, created_at
		FROM 
			posts WHERE id=$1`, id,
	).Scan(&p.ID, &p.Title, &p.Content, &p.AuthorID, &p.CreatedAt); err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *repo) DeleteByIDOwner(ctx context.Context, id, ownerID int64) (bool, error) {
	cmd, err := r.db.Pool.Exec(ctx, `
		DELETE FROM posts 
		WHERE id=$1 AND author_id=$2`, id, ownerID)
	return cmd.RowsAffected() > 0, err
}
