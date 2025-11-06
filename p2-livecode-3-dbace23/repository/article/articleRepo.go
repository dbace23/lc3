package articlerepo

import (
	"context"
	"instagram/model"
	"instagram/util/database"
)

type Repo interface {
	Create(ctx context.Context, a *model.Article) error
	All(ctx context.Context) ([]model.Article, error)
	ByID(ctx context.Context, id int64) (*model.Article, error)
	DeleteByIDOwner(ctx context.Context, id, ownerID int64) (bool, error)
}

type repo struct{ db *database.DB }

func New(db *database.DB) Repo { return &repo{db} }

func (r *repo) Create(ctx context.Context, a *model.Article) error {
	return r.db.Pool.QueryRow(ctx, `
		INSERT INTO 
		articles(title, content, author_id, category_id)
		VALUES 
			($1,$2,$3,$4) 
		RETURNING 
			id, created_at`,
		a.Title, a.Content, a.AuthorID, a.CategoryID,
	).Scan(&a.ID, &a.CreatedAt)
}

func (r *repo) All(ctx context.Context) ([]model.Article, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT 
			id, title, content, author_id, category_id, created_at
		FROM 
			articles 
		ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.Article
	for rows.Next() {
		var a model.Article
		if err := rows.Scan(&a.ID, &a.Title, &a.Content, &a.AuthorID, &a.CategoryID, &a.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (r *repo) ByID(ctx context.Context, id int64) (*model.Article, error) {
	var a model.Article
	if err := r.db.Pool.QueryRow(ctx, `
		SELECT 
			id, title, content, author_id, category_id, created_at
		FROM 
			articles
		WHERE id=$1`, id,
	).Scan(&a.ID, &a.Title, &a.Content, &a.AuthorID, &a.CategoryID, &a.CreatedAt); err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *repo) DeleteByIDOwner(ctx context.Context, id, ownerID int64) (bool, error) {
	cmd, err := r.db.Pool.Exec(ctx, `
		DELETE FROM articles 
		WHERE id=$1 AND author_id=$2`, id, ownerID)
	return cmd.RowsAffected() > 0, err
}
