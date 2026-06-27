package handler

import (
	"context"
	"errors"
	"image"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/traP-jp/h26s_02/domain"
	"github.com/traP-jp/h26s_02/repository"
	"github.com/traP-jp/h26s_02/storage"

	_ "image/jpeg" // フォーマット確認のために blank import
	_ "image/png"
)

type Post struct {
	db                 repository.DB
	postRepository     repository.Post
	reactionRepository repository.Reaction
	tagRepository      repository.Tag
	imageStorage       storage.Image
}

func NewPost(db repository.DB, postRepository repository.Post, reactionRepository repository.Reaction, tagRepository repository.Tag, imageStorage storage.Image) *Post {
	return &Post{
		db:                 db,
		postRepository:     postRepository,
		tagRepository:      tagRepository,
		reactionRepository: reactionRepository,
		imageStorage:       imageStorage,
	}
}

type PostPostResponse struct {
	ID uuid.UUID `json:"id"`
}

func (h *Post) GetPost(c *echo.Context) error {

	postID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid post ID")
	}
	post, err := h.postRepository.GetPost(c.Request().Context(), postID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid post ID")
	}
	reactions, err := h.reactionRepository.GetReactionCount(c.Request().Context(), postID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid post ID")
	}
	tags, err := h.tagRepository.GetPostTags(c.Request().Context(), postID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid post ID")
	}
	return c.JSON(http.StatusOK, struct {
		ID        uuid.UUID               `json:"id"`
		UserName  string                  `json:"user_name"`
		CreatedAt time.Time               `json:"created_at"`
		Reactions []*domain.ReactionCount `json:"reactions"`
		Tags      []string                `json:"tags"`
	}{
		ID:        post.GetID(),
		UserName:  post.GetUserName(),
		CreatedAt: post.GetCreatedAt(),
		Reactions: reactions,
		Tags:      tags,
	})
}

func (p *Post) PostPost(c *echo.Context) error {
	header, err := c.FormFile("image")
	if errors.Is(err, http.ErrMissingFile) {
		return echo.NewHTTPError(http.StatusBadRequest, "image is required")
	}
	if err != nil {
		log.Printf("failed to read form file: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	tagsString := c.FormValue("tags")
	var tags []string
	if len(tagsString) != 0 {
		tags = strings.Split(tagsString, ",")
		if len(tags) > 10 {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid tag count")
		}
		for _, tag := range tags {
			tagLength := len([]rune(tag))
			if tagLength == 0 || tagLength > 16 {
				return echo.NewHTTPError(http.StatusBadRequest, "invalid tag length")
			}
		}
	}

	f, err := header.Open()
	if err != nil {
		log.Printf("failed to open file: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	defer f.Close()

	_, format, err := image.DecodeConfig(f)
	if err != nil {
		log.Printf("failed to decode image config: %v\n", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid image format")
	}
	if format != "jpeg" && format != "png" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid image format")
	}
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		log.Printf("failed to rewind file: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	postID, err := uuid.NewV7()
	if err != nil {
		log.Printf("failed to generate post ID: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	userName, err := GetUserName(c)
	if err != nil {
		return err
	}

	err = p.db.Transaction(c.Request().Context(), func(ctx context.Context) error {
		err := p.postRepository.CreatePost(ctx, postID, userName)
		if err != nil {
			log.Printf("failed to create post: %v\n", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}

		err = p.tagRepository.CreatePostTags(ctx, postID, tags)
		if err != nil {
			log.Printf("failed to create post tags: %v\n", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}

		err = p.imageStorage.SaveImage(ctx, postID.String(), f)
		if err != nil {
			log.Printf("failed to save image: %v\n", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}

		return nil
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, PostPostResponse{
		ID: postID,
	})
}
