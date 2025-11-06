package userrepo

import (
	"context"
	"errors"

	"instagram/model"
	"instagram/util/database"
)

type Repo interface {
	Create(ctx context.Context, u *model.User) error
	ByEmail(ctx context.Context, email string) (*model.User, error)
}

type repo struct{ db *database.DB }

func New(db *database.DB) Repo { return &repo{db} }

func (r *repo) Create(ctx context.Context, u *model.User) error {
	return r.db.Pool.QueryRow(ctx, `
		INSERT INTO users(first_name, last_name, email, username, password_hash, age, address)
		VALUES ($1,$2,$3,$4,$5,18,'-') RETURNING id, created_at`,
		u.FirstName, u.LastName, u.Email, u.Username, u.PasswordHash,
	).Scan(&u.ID, &u.CreatedAt)
}

func (r *repo) ByEmail(ctx context.Context, email string) (*model.User, error) {
	var u model.User
	err := r.db.Pool.QueryRow(ctx, `
		SELECT 
			id, first_name, last_name, email, username, password_hash, created_at
		FROM 
			users 
		WHERE email=$1`, email,
	).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Username, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &u, nil
}
