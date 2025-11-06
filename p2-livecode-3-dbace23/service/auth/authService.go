package authsvc

import (
	"context"
	"errors"

	"instagram/model"
	userrepo "instagram/repository/user"
	"instagram/util/hash"
	jwtutil "instagram/util/jwt"

	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

var (
	ErrEmailTaken   = errors.New("email already registered")
	ErrBadInput     = errors.New("bad input")
	ErrInvalidCreds = errors.New("invalid credentials")
)

type Service interface {
	Register(ctx context.Context, req model.RegisterReq, secret string) (*model.User, string, error)
	Login(ctx context.Context, req model.LoginReq, secret string) (*model.User, string, error)
}

type service struct{ ur userrepo.Repo }

func New(ur userrepo.Repo) Service { return &service{ur} }

func (s *service) Register(ctx context.Context, req model.RegisterReq, secret string) (*model.User, string, error) {
	hashed, err := hash.HashPassword(req.Password)
	if err != nil {
		return nil, "", err
	}

	u := &model.User{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: hashed,
	}

	if err := s.ur.Create(ctx, u); err != nil {
		if isDuplicateErr(err) {
			return nil, "", ErrEmailTaken
		}
		return nil, "", err
	}

	token, err := jwtutil.Issue(secret, u.ID, "user", 24)
	if err != nil {
		return nil, "", err
	}

	return u, token, nil
}

func (s *service) Login(ctx context.Context, req model.LoginReq, secret string) (*model.User, string, error) {
	u, err := s.ur.ByEmail(ctx, req.Email)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	if !hash.Check(req.Password, u.PasswordHash) {
		return nil, "", errors.New("invalid credentials")
	}

	token, err := jwtutil.Issue(secret, u.ID, "user", 24)
	if err != nil {
		return nil, "", err
	}

	return u, token, nil
}

// ////
func isDuplicateErr(err error) bool {

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}
	// Postgres
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return true
	}
	//MySQL
	var myErr *mysql.MySQLError
	if errors.As(err, &myErr) && myErr.Number == 1062 {
		return true
	}
	return false
}
