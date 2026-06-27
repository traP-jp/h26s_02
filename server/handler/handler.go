package handler

import (
	"log"

	"github.com/labstack/echo/v5"
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

	api := e.Group("/api")
	api.GET("/users/me", h.user.GetMe)
	// api.GET("/posts", h.post.GetPosts)

	log.Fatal(e.Start(":8080"))
}
