package controller

import (
	"net/http"
	"strconv"

	"instagram/model"
	postsvc "instagram/service/post"

	"github.com/labstack/echo/v4"
)

type PostController struct {
	s postsvc.Service
}

func NewPostController(s postsvc.Service) *PostController {
	return &PostController{s}
}

// app/echoServer/controller/postController.go

// Create post
// @Summary      Create post
// @Security     BearerAuth
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        payload  body  model.CreatePostReq  true  "Create post payload"
// @Success      201  {object}  model.Post
// @Failure      400  {object}  map[string]any
// @Router       /v1/posts [post]
func (ct *PostController) Create(c echo.Context) error {
	claims := c.Get("claims").(map[string]any)
	uid := int64(claims["sub"].(float64))

	var req model.CreatePostReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid body"})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "validation error",
			"errors":  err.Error(),
		})
	}

	p, err := ct.s.Create(c.Request().Context(), uid, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
	}
	return c.JSON(http.StatusCreated, p)
}

// List posts
// @Summary      List posts
// @Security     BearerAuth
// @Tags         posts
// @Produce      json
// @Success      200  {array}   model.Post
// @Router       /v1/posts [get]
func (ct *PostController) List(c echo.Context) error {
	out, err := ct.s.List(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": err.Error()})
	}
	return c.JSON(http.StatusOK, out)
}

// Get post detail
// @Summary      Post detail
// @Security     BearerAuth
// @Tags         posts
// @Produce      json
// @Param        id   path  int  true  "Post ID"
// @Success      200  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /v1/posts/{id} [get]
func (ct *PostController) Detail(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	data, err := ct.s.Detail(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"message": "post not found"})
	}
	return c.JSON(http.StatusOK, data)
}

// Delete post
// @Summary      Delete post
// @Security     BearerAuth
// @Tags         posts
// @Produce      json
// @Param        id   path  int  true  "Post ID"
// @Success      200  {object}  map[string]any
// @Failure      403  {object}  map[string]any
// @Router       /v1/posts/{id} [delete]
func (ct *PostController) Delete(c echo.Context) error {
	claims := c.Get("claims").(map[string]any)
	uid := int64(claims["sub"].(float64))

	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := ct.s.Delete(c.Request().Context(), id, uid); err != nil {
		return c.JSON(http.StatusForbidden, echo.Map{"message": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "deleted", "id": id})
}
