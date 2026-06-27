package handler

import (
	"context"
	"errors"
	"image"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/traP-jp/h26s_02/repository"
	"github.com/traP-jp/h26s_02/storage"

	_ "image/jpeg"
	_ "image/png"
)

type Post struct {
	db             repository.DB
	postRepository repository.Post
	tagRepository  repository.Tag
	imageStorage   storage.Image
}

func NewPost(db repository.DB, postRepository repository.Post, tagRepository repository.Tag, imageStorage storage.Image) *Post {
	return &Post{
		db:             db,
		postRepository: postRepository,
		tagRepository:  tagRepository,
		imageStorage:   imageStorage,
	}
}

type PostPostResponse struct {
	ID uuid.UUID `json:"id"`
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
