package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/traP-jp/h26s_02/repository"
)

type Post struct {
	postRepository repository.Post
}

func NewPost(postRepository repository.Post) *Post {
	return &Post{
		postRepository: postRepository,
	}
}
func (h *Post) GetPost(c *echo.Context) error {

	postID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid post ID")
	}
	post, err := h.postRepository.GetPost(postID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid post ID")
	}
	return c.JSON(http.StatusOK, post)
}
