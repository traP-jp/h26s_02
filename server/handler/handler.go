package handler

import (
	"log"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type Handler struct {
	user *User
	post *Post
}

func NewHandler(user *User, post *Post) *Handler {
	return &Handler{
		user: user,
		post: post,
	}
}

func (h *Handler) Start(e *echo.Echo) {
	e.Use(middleware.Recover(), middleware.RequestLogger())

	api := e.Group("/api", AuthMiddleware())
	api.GET("/users/me", h.user.GetMe)
	api.POST("/posts", h.post.PostPost)
	// api.GET("/posts", h.post.GetPosts)

	log.Fatal(e.Start(":8080"))
}
