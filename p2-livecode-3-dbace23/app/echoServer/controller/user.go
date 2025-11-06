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
func (ct *UserController) Register(c echo.Context) error {
	var req model.RegisterReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid body"})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "validation error",
			"errors":  err.Error(),
		})
	}
	u, token, err := ct.s.Register(c.Request().Context(), req, ct.jwtSecret)
	if err != nil {
		switch {
		case errors.Is(err, authsvc.ErrEmailTaken):
			return c.JSON(http.StatusConflict, echo.Map{"message": "email already registered"})
		case errors.Is(err, authsvc.ErrBadInput):
			return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
		default:
			rid := c.Response().Header().Get(echo.HeaderXRequestID)
			ct.log.Error("register failed",
				slog.Any("err", err),
				slog.String("req_id", rid),
				slog.String("path", c.Path()),
				slog.String("method", c.Request().Method),
			)
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "internal server error"})
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
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid body"})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "validation error",
			"errors":  err.Error(),
		})
	}
	u, token, err := ct.s.Login(c.Request().Context(), req, ct.jwtSecret)
	if err != nil {
		switch {
		case errors.Is(err, authsvc.ErrInvalidCreds):
			return c.JSON(http.StatusUnauthorized, echo.Map{"message": "invalid email or password"})
		case errors.Is(err, authsvc.ErrBadInput):
			return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
		default:
			rid := c.Response().Header().Get(echo.HeaderXRequestID)
			ct.log.Error("login failed",
				slog.Any("err", err),
				slog.String("req_id", rid),
				slog.String("path", c.Path()),
				slog.String("method", c.Request().Method),
			)
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "internal server error"})
		}
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "login success",
		"user":    u,
		"token":   token,
	})
}
