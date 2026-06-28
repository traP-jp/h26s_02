package handler

import (
	"context"
	"errors"
	"image"
	"io"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/traP-jp/h26s_02/repository"
	"github.com/traP-jp/h26s_02/storage"

	_ "image/jpeg" // フォーマット確認のために blank import
	_ "image/png"
)

type reactionCountResponse struct {
	ID    int `json:"id"`
	Count int `json:"count"`
}
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

func (p *Post) GetPost(c *echo.Context) error {

	postID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid post ID")
	}
	post, err := p.postRepository.GetPost(c.Request().Context(), postID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid post ID")
	}
	reactions, err := p.reactionRepository.GetReactionCount(c.Request().Context(), postID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid post ID")
	}

	userName, err := GetUserName(c)
	if err != nil {
		return err
	}
	userReactions, err := p.reactionRepository.GetUserReactionsByPostIDs(c.Request().Context(), userName, []uuid.UUID{postID})
	if err != nil {
		log.Printf("failed to get user reactions: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	postUserReactions, ok := userReactions[postID]
	if !ok {
		postUserReactions = []int{}
	}

	reactionCountResponses := make([]ReactionResponse, 0, len(reactions))
	for _, reaction := range reactions {
		myReaction := slices.Contains(postUserReactions, reaction.GetID())
		reactionCountResponses = append(reactionCountResponses, ReactionResponse{
			ID:         reaction.GetID(),
			Count:      reaction.GetCount(),
			MyReaction: myReaction,
		})
	}
	tags, err := p.tagRepository.GetPostTags(c.Request().Context(), postID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid post ID")
	}

	imageURL, err := p.imageStorage.GetTemporalyURL(c.Request().Context(), postID.String())
	if errors.Is(err, storage.ErrImageNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "image not found")
	}
	if err != nil {
		log.Printf("failed to get iamge temporaly url: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusOK, PostResponse{
		ID:        post.GetID(),
		UserName:  post.GetUserName(),
		CreatedAt: post.GetCreatedAt(),
		Reactions: reactionCountResponses,
		Tags:      tags,
		ImageURL:  imageURL,
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

func (p *Post) DeleteReaction(c *echo.Context) error {
	postID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid post ID")
	}

	userName, err := GetUserName(c)
	if err != nil {
		return err
	}
	deleteReactionID, err := strconv.Atoi(c.Param("reaction_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid reaction ID")
	}
	err = p.reactionRepository.DeleteReaction(c.Request().Context(), postID, userName, deleteReactionID)
	if errors.Is(err, repository.ErrNoRecordDeleted) {
		return echo.NewHTTPError(http.StatusNotFound, "reaction not found")
	} else if err != nil {
		log.Printf("failed to delete reaction: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	return c.NoContent(http.StatusOK)
}

type GetPostsRequest struct {
	Before uuid.UUID `query:"before"`
	Limit  int       `query:"limit"`
}

type ReactionResponse struct {
	ID         int  `json:"id"`
	Count      int  `json:"count"`
	MyReaction bool `json:"myReaction"`
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
	ctx := c.Request().Context()

	req := GetPostsRequest{Limit: 30}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid query parameters")
	}

	var referenceTime time.Time
	if req.Before != uuid.Nil {
		post, err := p.postRepository.GetPostByID(ctx, req.Before)
		if err != nil {
			if errors.Is(err, repository.ErrRecordNotFound) {
				return echo.NewHTTPError(http.StatusNotFound, "post not found")
			}
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}
		referenceTime = post.GetCreatedAt()
	} else {
		referenceTime = time.Now()
	}

	posts, err := p.postRepository.GetPosts(ctx, referenceTime, req.Limit)
	if err != nil {
		log.Printf("failed to get posts: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	if len(posts) == 0 {
		return c.JSON(http.StatusOK, []PostResponse{})
	}

	postIDs := make([]uuid.UUID, len(posts))
	for i, post := range posts {
		postIDs[i] = post.GetID()
	}

	allTags, err := p.tagRepository.GetTagsByPostIDs(ctx, postIDs)
	if err != nil {
		log.Printf("failed to get tags: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	allReactions, err := p.reactionRepository.GetReactionsByPostIDs(ctx, postIDs)
	if err != nil {
		log.Printf("failed to get reactions: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	userName, err := GetUserName(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}
	userReactions, err := p.reactionRepository.GetUserReactionsByPostIDs(ctx, userName, postIDs)
	if err != nil {
		log.Printf("failed to get user reactions: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	response := make([]PostResponse, 0, len(posts))
	for _, post := range posts {
		postID := post.GetID()

		var tagNames []string
		for _, t := range allTags[postID] {
			tagNames = append(tagNames, t.GetName())
		}

		var reactionRes []ReactionResponse
		for _, r := range allReactions[postID] {
			if r.GetCount() > 0 {
				myReaction := slices.Contains(userReactions[postID], r.GetID())
				reactionRes = append(reactionRes, ReactionResponse{
					ID:         r.GetID(),
					Count:      r.GetCount(),
					MyReaction: myReaction,
				})
			}
		}

		imageURL, err := p.imageStorage.GetTemporalyURL(ctx, postID.String())
		if errors.Is(err, storage.ErrImageNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "image not found")
		}
		if err != nil {
			log.Printf("failed to get iamge temporaly url: %v\n", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}

		response = append(response, PostResponse{
			ID:        postID,
			UserName:  post.GetUserName(),
			Tags:      tagNames,
			ImageURL:  imageURL,
			Reactions: reactionRes,
			CreatedAt: post.GetCreatedAt(),
		})
	}

	return c.JSON(http.StatusOK, response)
}
func (p *Post) GetPostsByUser(c *echo.Context) error {
	userName := c.Param("id")
	ctx := c.Request().Context()

	posts, err := p.postRepository.GetPostsByUser(ctx, userName)
	if err != nil {
		log.Printf("failed to get posts: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	if len(posts) == 0 {
		return c.JSON(http.StatusOK, []PostResponse{})
	}

	postIDs := make([]uuid.UUID, len(posts))
	for i, post := range posts {
		postIDs[i] = post.GetID()
	}

	allTags, err := p.tagRepository.GetTagsByPostIDs(ctx, postIDs)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get tags")
	}
	allReactions, err := p.reactionRepository.GetReactionsByPostIDs(ctx, postIDs)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get reactions")
	}

	loginUserName, err := GetUserName(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}
	loginUserReactions, err := p.reactionRepository.GetUserReactionsByPostIDs(ctx, loginUserName, postIDs)
	if err != nil {
		log.Printf("failed to get user reactions: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	response := make([]PostResponse, 0, len(posts))
	for _, post := range posts {
		postID := post.GetID()

		var tagNames []string
		for _, t := range allTags[postID] {
			tagNames = append(tagNames, t.GetName())
		}

		var reactionRes []ReactionResponse
		for _, r := range allReactions[postID] {
			if r.GetCount() > 0 {
				myReaction := slices.Contains(loginUserReactions[postID], r.GetID())
				reactionRes = append(reactionRes, ReactionResponse{
					ID:         r.GetID(),
					Count:      r.GetCount(),
					MyReaction: myReaction,
				})
			}
		}

		imageURL, err := p.imageStorage.GetTemporalyURL(ctx, postID.String())
		if errors.Is(err, storage.ErrImageNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "image not found")
		}
		if err != nil {
			log.Printf("failed to get iamge temporaly url: %v\n", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}

		response = append(response, PostResponse{
			ID:        postID,
			UserName:  post.GetUserName(),
			Tags:      tagNames,
			ImageURL:  imageURL,
			Reactions: reactionRes,
			CreatedAt: post.GetCreatedAt(),
		})
	}

	return c.JSON(http.StatusOK, response)

}
