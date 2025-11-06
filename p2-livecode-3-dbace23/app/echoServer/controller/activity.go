package controller

import (
	activitysvc "instagram/service/activity"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ActivityController struct{ s activitysvc.Service }

func NewActivityController(s activitysvc.Service) *ActivityController { return &ActivityController{s} }

// app/echoServer/controller/activityController.go

// List my activities
// @Summary      My activities
// @Security     BearerAuth
// @Tags         activities
// @Produce      json
// @Success      200  {array}  model.Activity
// @Router       /v1/activities [get]
func (ct *ActivityController) ListMine(c echo.Context) error {
	uid := c.Get("uid").(int64)
	acts, err := ct.s.ListMine(c.Request().Context(), uid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
	}
	return c.JSON(http.StatusOK, acts)
}
