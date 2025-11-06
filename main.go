// Package main Instagram-like API.
//
// @title           Instagram Mini API
// @version         1.0
// @description     Small Instagram-like service (posts, likes, activities, users).
// @contact.name    Halim Iskandar
// @contact.email   halim.iskandar2323@gmail.com
// @BasePath        /
// @schemes         http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description  Use:  Bearer <JWT>
package main

import (
	"context"
	echoServer "instagram/app/echoServer"
	"instagram/app/echoServer/controller"
	"instagram/app/echoServer/validation"
	"instagram/config"
	activityrepo "instagram/repository/activity"
	jokerepo "instagram/repository/joke"
	likerepo "instagram/repository/like"
	postrepo "instagram/repository/post"
	userrepo "instagram/repository/user"
	activitysvc "instagram/service/activity"
	authsvc "instagram/service/auth"
	likesvc "instagram/service/like"
	postsvc "instagram/service/post"
	"instagram/util/database"
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {
	secret := os.Getenv("JWT_SECRET")
	cfg := config.Load()
	ctx := context.Background()

	db, err := database.New(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("db connect failed", "err", err)
		os.Exit(1)
	}
	defer db.Pool.Close()

	// repos
	pr := postrepo.New(db)
	lr := likerepo.New(db)
	ar := activityrepo.New(db)
	ur := userrepo.New(db)
	jr := jokerepo.New(cfg.ApiNinjasKey)

	// services
	ps := postsvc.New(pr, lr, ar, jr)
	ls := likesvc.New(lr, pr, ar)
	as := activitysvc.New(ar)
	aus := authsvc.New(ur)

	slog.Info("JWT secret check",
    "secret_length", len(cfg.JWTSecret),
    "secret_preview", cfg.JWTSecret[:min(6, len(cfg.JWTSecret))],
)
	// controllers
	pc := controller.NewPostController(ps)
	lc := controller.NewLikeController(ls)
	ac := controller.NewActivityController(as)
	uc := controller.NewUserController(aus, secret, slog.Default())

	// echo
	e := echo.New()
	echoServer.RegisterMiddlewares(e)
	e.Validator = validation.New()

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]any{
			"status":  "ok",
			"message": "✅ Service is healthy and connected",
		})
	})

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	echoServer.Register(e, echoServer.C{
		User:      uc,
		Post:      pc,
		Like:      lc,
		Activity:  ac,
		JWTSecret: secret,
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Port
	}
	if port == "" {
		port = "8080"
	}

	slog.Info("starting server", "PORT_env", os.Getenv("PORT"), "chosen_port", port)

	e.Logger.Fatal(e.Start(":" + port))
}
