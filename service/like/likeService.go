package likesvc

import (
	"context"
	"errors"
	"fmt"

	"instagram/model"
	activityrepo "instagram/repository/activity"
	likerepo "instagram/repository/like"
	postrepo "instagram/repository/post"
)

type Service interface {
	Create(ctx context.Context, userID int64, req model.CreateLikeReq) (*model.Like, error)
	Detail(ctx context.Context, id int64) (*model.Like, error)
	Delete(ctx context.Context, id, userID int64) error
}

type service struct {
	lr  likerepo.Repo
	pr  postrepo.Repo
	log activityrepo.Repo
}

func New(lr likerepo.Repo, pr postrepo.Repo, log activityrepo.Repo) Service {
	return &service{lr, pr, log}
}

func (s *service) Create(ctx context.Context, userID int64, req model.CreateLikeReq) (*model.Like, error) {
	if _, err := s.pr.ByID(ctx, req.PostID); err != nil {
		return nil, errors.New("post not found")
	}
	lk, err := s.lr.Create(ctx, userID, req.PostID)
	if err != nil {
		return nil, err
	}
	_ = s.log.Log(ctx, model.Activity{UserID: userID, Action: "LIKE_CREATE", Description: fmt.Sprintf("like POST id=%d", req.PostID)})
	return lk, nil
}
func (s *service) Detail(ctx context.Context, id int64) (*model.Like, error) {
	return s.lr.ByID(ctx, id)
}
func (s *service) Delete(ctx context.Context, id, userID int64) error {
	ok, err := s.lr.DeleteByIDOwner(ctx, id, userID)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("not allowed or like not found")
	}
	_ = s.log.Log(ctx, model.Activity{UserID: userID, Action: "LIKE_DELETE", Description: fmt.Sprintf("unlike like_id=%d", id)})
	return nil
}
