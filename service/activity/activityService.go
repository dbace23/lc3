package activitysvc

import (
	"context"

	"instagram/model"
	activityrepo "instagram/repository/activity"
)

type Service interface {
	ListMine(ctx context.Context, userID int64) ([]model.Activity, error)
}

type service struct{ ar activityrepo.Repo }

func New(ar activityrepo.Repo) Service { return &service{ar} }

func (s *service) ListMine(ctx context.Context, userID int64) ([]model.Activity, error) {
	return s.ar.ListByUser(ctx, userID)
}
