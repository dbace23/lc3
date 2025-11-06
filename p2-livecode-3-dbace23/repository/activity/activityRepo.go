package activityrepo

import (
	"context"

	"instagram/model"
	"instagram/util/database"
)

type Repo interface {
	Log(ctx context.Context, a model.Activity) error
	ListByUser(ctx context.Context, userID int64) ([]model.Activity, error)
}

type repo struct{ db *database.DB }

func New(db *database.DB) Repo { return &repo{db} }

func (r *repo) Log(ctx context.Context, a model.Activity) error {
	_, err := r.db.Pool.Exec(ctx, `
		INSERT INTO user_activity_logs(user_id, action, description)
		VALUES ($1,$2,$3)`, a.UserID, a.Action, a.Description)
	return err
}

func (r *repo) ListByUser(ctx context.Context, userID int64) ([]model.Activity, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT 
			id, user_id, action, description, created_at
		FROM 
			user_activity_logs
		WHERE 
			user_id=$1 
		ORDER BY id DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.Activity
	for rows.Next() {
		var a model.Activity
		if err := rows.Scan(&a.ID, &a.UserID, &a.Action, &a.Description, &a.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}
