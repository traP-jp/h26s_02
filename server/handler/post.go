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

type GetPostsRequest struct {
	Before uuid.UUID `query:"before"`
	Limit  int       `query:"limit"`
}

type ReactionResponse struct {
	ID    int `json:"id"`
	Count int `json:"count"`
}

type PostResponse struct {
	ID        uuid.UUID          `json:"id"`
	UserName  string             `json:"userName"`
	Tags      []string           `json:"tags"`
	ImageURL  string             `json:"imageUrl"`
	Reactions []ReactionResponse `json:"reactions"`
	CreatedAt time.Time          `json:"createdAt"`
}

func (p *Post) GetPosts(c *echo.Context) error {
	req := GetPostsRequest{Limit: 30}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid query parameters")
	}

	var referenceTime time.Time
	if req.Before != uuid.Nil {
		post, err := p.postRepository.GetPostByID(req.Before.String())
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "post not found")
		}
		referenceTime = post.GetCreatedAt()
	} else {
		referenceTime = time.Now()
	}

	posts, err := p.postRepository.GetPosts(referenceTime, req.Limit)
	if err != nil {
		log.Printf("failed to get posts: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	response := make([]PostResponse, 0, len(posts))
	for _, post := range posts {
		reactions := make([]ReactionResponse, 0)
		for _, reaction := range post.GetReactions() {
			if reaction.GetCount() > 0 {
				reactions = append(reactions, ReactionResponse{
					ID:    reaction.GetID(),
					Count: reaction.GetCount(),
				})
			}
		}

		response = append(response, PostResponse{
			ID:        post.GetID(),
			UserName:  post.GetUserName(),
			Tags:      post.GetTags(),
			ImageURL:  "",
			Reactions: reactions,
			CreatedAt: post.GetCreatedAt(),
		})
	}

	return c.JSON(http.StatusOK, response)
}
