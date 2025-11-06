// app/echoServer/controller/postController.go
package controller

import (
	"errors"
	"net/http"
	"strconv"

	"instagram/model"
	postsvc "instagram/service/post"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type PostController struct{ s postsvc.Service }

func NewPostController(s postsvc.Service) *PostController { return &PostController{s} }

// Create post
// @Summary      Create post
// @Description  Create a new post (JWT required)
// @Security     BearerAuth
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        payload  body  model.CreatePostReq  true  "Create post payload"
// @Success      201  {object}  model.Post
// @Failure      400  {object}  map[string]any "validation error / bad input"
// @Failure      401  {object}  map[string]any "missing or invalid token"
// @Failure      500  {object}  map[string]any "internal server error"
// @Router       /v1/posts [post]
func (ct *PostController) Create(c echo.Context) error {

	tok, ok := c.Get("user").(*jwt.Token)
	if !ok || tok == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "missing or invalid token")
	}
	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid token claims")
	}
	sub, ok := claims["sub"].(float64)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid subject in token")
	}
	userID := int64(sub)

	var req model.CreatePostReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}
	if err := c.Validate(&req); err != nil {

		return echo.NewHTTPError(http.StatusBadRequest, "validation error")
	}

	p, err := ct.s.Create(c.Request().Context(), userID, req)
	if err != nil {
		switch {
		case errors.Is(err, postsvc.ErrBadInput):
			return echo.NewHTTPError(http.StatusBadRequest, postsvc.ErrBadInput.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusCreated, p)
}

// List posts
// @Summary      List posts
// @Description  List all posts (JWT required)
// @Security     BearerAuth
// @Tags         posts
// @Produce      json
// @Success      200  {array}   model.Post
// @Failure      401  {object}  map[string]any "missing or invalid token"
// @Failure      500  {object}  map[string]any "internal server error"
// @Router       /v1/posts [get]
func (ct *PostController) List(c echo.Context) error {
	out, err := ct.s.List(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, out)
}

// Get post detail
// @Summary      Post detail
// @Description  Get a post by ID (JWT required)
// @Security     BearerAuth
// @Tags         posts
// @Produce      json
// @Param        id   path  int  true  "Post ID"
// @Success      200  {object}  map[string]any
// @Failure      400  {object}  map[string]any "invalid id"
// @Failure      401  {object}  map[string]any "missing or invalid token"
// @Failure      404  {object}  map[string]any "post not found"
// @Failure      500  {object}  map[string]any "internal server error"
// @Router       /v1/posts/{id} [get]
func (ct *PostController) Detail(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}
	data, err := ct.s.Detail(c.Request().Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, postsvc.ErrNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "post not found")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusOK, data)
}

// Delete post
// @Summary      Delete post
// @Description  Delete a post by ID (JWT required; only owner can delete)
// @Security     BearerAuth
// @Tags         posts
// @Produce      json
// @Param        id   path  int  true  "Post ID"
// @Success      200  {object}  map[string]any "deleted"
// @Failure      400  {object}  map[string]any "invalid id"
// @Failure      401  {object}  map[string]any "missing or invalid token"
// @Failure      403  {object}  map[string]any "forbidden - not owner"
// @Failure      404  {object}  map[string]any "post not found"
// @Failure      500  {object}  map[string]any "internal server error"
// @Router       /v1/posts/{id} [delete]
func (ct *PostController) Delete(c echo.Context) error {
	tok, ok := c.Get("user").(*jwt.Token)
	if !ok || tok == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "missing or invalid token")
	}
	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid token claims")
	}
	sub, ok := claims["sub"].(float64)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid subject in token")
	}
	userID := int64(sub)

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}
	if err := ct.s.Delete(c.Request().Context(), id, userID); err != nil {
		switch {
		case errors.Is(err, postsvc.ErrNotOwner):
			return echo.NewHTTPError(http.StatusForbidden, "forbidden: not owner")
		case errors.Is(err, postsvc.ErrNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "post not found")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "deleted", "id": id})
}
