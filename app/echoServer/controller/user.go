package controller

import (
	"errors"
	"log/slog"
	"net/http"

	"instagram/model"
	authsvc "instagram/service/auth"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	s         authsvc.Service
	jwtSecret string
	log       *slog.Logger
}

func NewUserController(s authsvc.Service, secret string, log *slog.Logger) *UserController {
	return &UserController{
		s:         s,
		jwtSecret: secret,
		log:       log,
	}
}

// app/echoServer/controller/userController.go

// Register a new user
// @Summary      Register user
// @Description  Register a new user with email validation
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        payload  body  model.RegisterReq  true  "Register payload"
// @Success      201  {object}  map[string]any
// @Failure      400  {object}  map[string]any
// @Failure      409 {object} map[string]any "email already registered"
// @Failure      500 {object} map[string]any "internal server error"
// @Router       /v1/users/register [post]
// Register a new user
func (ct *UserController) Register(c echo.Context) error {
	var req model.RegisterReq
	if err := c.Bind(&req); err != nil {
		ct.log.Warn("bind failed", "path", c.Path(), "err", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}
	if err := c.Validate(&req); err != nil {
		ct.log.Warn("validation failed", "path", c.Path(), "err", err)

		return echo.NewHTTPError(http.StatusBadRequest)
	}

	u, token, err := ct.s.Register(c.Request().Context(), req, ct.jwtSecret)
	if err != nil {
		switch {
		case errors.Is(err, authsvc.ErrEmailTaken):
			return echo.NewHTTPError(http.StatusConflict, "email already registered")
		case errors.Is(err, authsvc.ErrBadInput):

			ct.log.Warn("bad input", "path", c.Path(), "err", err)
			return echo.NewHTTPError(http.StatusBadRequest)
		default:
			rid := c.Response().Header().Get(echo.HeaderXRequestID)
			ct.log.Error("register failed", "err", err, "req_id", rid, "path", c.Path(), "method", c.Request().Method)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "registered",
		"user":    u,
		"token":   token,
	})
}

// Login
// @Summary      Login
// @Description  Login with email + password, returns JWT
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        payload  body  model.LoginReq  true  "Login payload"
// @Success      200  {object}  map[string]any
// @Failure      400  {object}  map[string]any
// @Failure      401 {object} map[string]any
// @Failure      500 {object} map[string]any
// @Router       /v1/users/login [post]
func (ct *UserController) Login(c echo.Context) error {
	var req model.LoginReq
	if err := c.Bind(&req); err != nil {
		ct.log.Warn("bind failed", "path", c.Path(), "err", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}
	if err := c.Validate(&req); err != nil {
		ct.log.Warn("validation failed", "path", c.Path(), "err", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	u, token, err := ct.s.Login(c.Request().Context(), req, ct.jwtSecret)
	if err != nil {
		switch {
		case errors.Is(err, authsvc.ErrInvalidCreds):
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid email or password")
		case errors.Is(err, authsvc.ErrBadInput):
			ct.log.Warn("bad input", "path", c.Path(), "err", err)
			return echo.NewHTTPError(http.StatusBadRequest)
		default:
			rid := c.Response().Header().Get(echo.HeaderXRequestID)
			ct.log.Error("login failed", "err", err, "req_id", rid, "path", c.Path(), "method", c.Request().Method)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "login success",
		"user":    u,
		"token":   token,
	})
}
