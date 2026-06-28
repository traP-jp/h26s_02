package handler

import (
	"log"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type Handler struct {
	user *User
	post *Post
	tag  *Tag
}

func NewHandler(user *User, post *Post, tag *Tag) *Handler {
	return &Handler{
		user: user,
		post: post,
		tag:  tag,
	}
}

func (h *Handler) Start(e *echo.Echo) {
	e.Use(middleware.Recover(), middleware.RequestLogger())

	api := e.Group("/api", AuthMiddleware())
	api.GET("/users/me", h.user.GetMe)
	api.POST("/posts", h.post.PostPost)
	api.GET("/posts/:id", h.post.GetPost)
	api.POST("/posts/:id/reactions", h.post.PostReaction)
	api.GET("/tags", h.tag.GetTags)

	api.DELETE("/posts/:id/reactions/:reaction_id", h.post.DeleteReaction)

	api.GET("/users/:id/posts", h.post.GetPostsByUser)

	// api.GET("/posts", h.post.GetPosts)

	api.GET("/posts", h.post.GetPosts)

	log.Fatal(e.Start(":8080"))
}
