package controller

// controller/likeController.go
import (
	"net/http"
	"strconv"

	"instagram/model"
	likesvc "instagram/service/like"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type LikeController struct{ s likesvc.Service }

func NewLikeController(s likesvc.Service) *LikeController { return &LikeController{s} }

func userIDFromJWT(c echo.Context) (int64, error) {
	tok, ok := c.Get("user").(*jwt.Token)
	if !ok || tok == nil {
		return 0, echo.NewHTTPError(http.StatusUnauthorized, "missing or invalid token")
	}
	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return 0, echo.NewHTTPError(http.StatusUnauthorized, "invalid token claims")
	}
	sub, ok := claims["sub"].(float64)
	if !ok {
		return 0, echo.NewHTTPError(http.StatusUnauthorized, "invalid subject in token")
	}
	return int64(sub), nil
}

func (ct *LikeController) Create(c echo.Context) error {
	uid, err := userIDFromJWT(c)
	if err != nil {
		return err
	}

	var req model.CreateLikeReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid body"})
	}

	if req.PostID == 0 {
		if idStr := c.Param("id"); idStr != "" {
			if v, err := strconv.ParseInt(idStr, 10, 64); err == nil {
				req.PostID = v
			}
		}
	}

	lk, err := ct.s.Create(c.Request().Context(), uid, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
	}
	return c.JSON(http.StatusCreated, echo.Map{"message": "liked", "data": lk})
}

func (ct *LikeController) Detail(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid id"})
	}
	lk, err := ct.s.Detail(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"message": "like not found"})
	}
	return c.JSON(http.StatusOK, lk)
}

func (ct *LikeController) Delete(c echo.Context) error {
	uid, err := userIDFromJWT(c)
	if err != nil {
		return err
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid id"})
	}
	if err := ct.s.Delete(c.Request().Context(), id, uid); err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"message": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "unliked", "id": id})
}
