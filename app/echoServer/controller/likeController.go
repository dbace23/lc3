package controller

import (
	"net/http"
	"strconv"

	"instagram/model"
	likesvc "instagram/service/like"

	"github.com/labstack/echo/v4"
)

type LikeController struct{ s likesvc.Service }

func NewLikeController(s likesvc.Service) *LikeController { return &LikeController{s} }

func (ct *LikeController) Create(c echo.Context) error {
	uid := c.Get("uid").(int64)
	var req model.CreateLikeReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid body"})
	}
	lk, err := ct.s.Create(c.Request().Context(), uid, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
	}
	return c.JSON(http.StatusCreated, echo.Map{"message": "liked", "data": lk})
}
func (ct *LikeController) Detail(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	lk, err := ct.s.Detail(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"message": "like not found"})
	}
	return c.JSON(http.StatusOK, lk)
}

func (ct *LikeController) Delete(c echo.Context) error {
	uid := c.Get("uid").(int64)
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := ct.s.Delete(c.Request().Context(), id, uid); err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"message": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "unliked", "id": id})
}
