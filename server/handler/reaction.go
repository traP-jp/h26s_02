package handler

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/traP-jp/h26s_02/repository"
)

type postReactionRequest struct {
	ID int `json:"id"`
}

func (h *Post) PostReaction(c *echo.Context) error {
	postID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid post ID")
	}

	var req postReactionRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	if req.ID < 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid reaction ID")
	}

	userName, err := GetUserName(c)
	if err != nil {
		return err
	}

	var reactionCountResponses []reactionCountResponse
	err = h.db.Transaction(c.Request().Context(), func(ctx context.Context) error {
		_, err := h.postRepository.GetPost(ctx, postID)
		if errors.Is(err, repository.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "post not found")
		}
		if err != nil {
			log.Printf("failed to get post: %v\n", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}

		if err := h.reactionRepository.CreateReaction(ctx, postID, req.ID, userName); err != nil {
			if errors.Is(err, repository.ErrUniqueKeyDuplicated) {
				return echo.NewHTTPError(http.StatusBadRequest, "already reacted")
			}
			if errors.Is(err, repository.ErrViolatedForeignKey) {
				return echo.NewHTTPError(http.StatusNotFound, "post not found")
			}
			log.Printf("failed to create reaction: %v\n", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}

		reactions, err := h.reactionRepository.GetReactionCount(ctx, postID)
		if err != nil {
			log.Printf("failed to get reaction count: %v\n", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}

		reactionCountResponses = make([]reactionCountResponse, 0, len(reactions))
		for _, reaction := range reactions {
			reactionCountResponses = append(reactionCountResponses, reactionCountResponse{
				ID:    reaction.GetID(),
				Count: reaction.GetCount(),
			})
		}

		return nil
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, reactionCountResponses)
}
