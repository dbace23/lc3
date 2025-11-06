// service/post/postService.go
package postsvc

import (
	"context"
	"fmt"
	"log/slog"

	"instagram/model"
	activityrepo "instagram/repository/activity"
	jokerrepo "instagram/repository/joke"
	likerepo "instagram/repository/like"
	postrepo "instagram/repository/post"
)

type Service interface {
	Create(ctx context.Context, userID int64, req model.CreatePostReq) (*model.Post, error)
	List(ctx context.Context) ([]model.Post, error)
	Detail(ctx context.Context, id int64) (map[string]any, error)
	Delete(ctx context.Context, id, userID int64) error
}

type service struct {
	pr       postrepo.Repo
	lr       likerepo.Repo
	log      activityrepo.Repo
	jokeRepo jokerrepo.Repo
}

func New(pr postrepo.Repo, lr likerepo.Repo, log activityrepo.Repo, jr jokerrepo.Repo) Service {
	return &service{pr: pr, lr: lr, log: log, jokeRepo: jr}
}

func (s *service) Create(ctx context.Context, userID int64, req model.CreatePostReq) (*model.Post, error) {
	if req.Title == "" {
		return nil, ErrBadInput
	}

	var content string
	if req.Content != nil {
		content = *req.Content
	}

	if s.jokeRepo != nil {
		if joke, err := s.jokeRepo.FetchJoke(ctx); err == nil && joke != "" {
			if content != "" {
				content += "\n\n" + "💡 Joke of the day: " + joke
			} else {
				content = "💡 Joke of the day: " + joke
			}
		} else if err != nil {
			slog.Warn("fetch joke failed", "err", err)
		}
	}

	p := &model.Post{
		Title:    req.Title,
		Content:  content,
		AuthorID: userID,
	}
	if err := s.pr.Create(ctx, p); err != nil {
		return nil, err
	}

	_ = s.log.Log(ctx, model.Activity{
		UserID:      userID,
		Action:      "POST_CREATE",
		Description: fmt.Sprintf("create POST id=%d title=%q", p.ID, p.Title),
	})
	return p, nil
}

func (s *service) List(ctx context.Context) ([]model.Post, error) {
	return s.pr.All(ctx)
}

func (s *service) Detail(ctx context.Context, id int64) (map[string]any, error) {
	post, err := s.pr.ByID(ctx, id)
	if err != nil || post == nil {
		return nil, ErrNotFound
	}

	likes, _ := s.lr.ListByPost(ctx, id)
	count, _ := s.lr.CountByPost(ctx, id)

	return map[string]any{
		"post":        post,
		"likes":       likes,
		"likes_count": count,
	}, nil
}

func (s *service) Delete(ctx context.Context, id, userID int64) error {

	ok, err := s.pr.DeleteByIDOwner(ctx, id, userID)
	if err != nil {
		return err
	}
	if ok {
		_ = s.log.Log(ctx, model.Activity{
			UserID:      userID,
			Action:      "POST_DELETE",
			Description: fmt.Sprintf("delete POST id=%d", id),
		})
		return nil
	}

	if p, err := s.pr.ByID(ctx, id); err == nil && p != nil {
		return ErrNotOwner
	}
	return ErrNotFound
}
