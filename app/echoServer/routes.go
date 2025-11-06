package echoServer

import (
	"instagram/app/echoServer/controller"

	"github.com/labstack/echo/v4"
)

type C struct {
	User     *controller.UserController
	Post     *controller.PostController
	Like     *controller.LikeController
	Activity *controller.ActivityController

	JWTSecret string
}

func Register(e *echo.Echo, c C) {
	pub := e.Group("v1")
	// --- User routes ---
	pub.POST("/users/register", c.User.Register)
	pub.POST("/users/login", c.User.Login)

	// --- Post routes ---
	pub.POST("/posts", c.Post.Create)
	pub.GET("/posts", c.Post.List)
	pub.GET("/posts/:id", c.Post.Detail)
	pub.DELETE("/posts/:id", c.Post.Delete)

	// --- Like routes ---
	pub.POST("/likes", c.Like.Create)
	pub.DELETE("/likes/:id", c.Like.Delete)

	// --- Activity routes ---
	pub.GET("/activities", c.Activity.ListMine)
}
