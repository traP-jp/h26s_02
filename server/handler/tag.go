package handler

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/traP-jp/h26s_02/repository"
)

type Tag struct {
	tagRepository repository.Tag
}

func NewTag(tagRepository repository.Tag) *Tag {
	return &Tag{
		tagRepository: tagRepository,
	}
}

type GetTagsResponse struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func (t *Tag) GetTags(c *echo.Context) error {
	tagCounts, err := t.tagRepository.GetTags(c.Request().Context())
	if err != nil {
		log.Printf("failed to get tags: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	response := make([]GetTagsResponse, 0, len(tagCounts))
	for _, tagCount := range tagCounts {
		response = append(response, GetTagsResponse{
			Name:  tagCount.Tag.GetName(),
			Count: tagCount.Count,
		})
	}

	return c.JSON(http.StatusOK, response)
}
