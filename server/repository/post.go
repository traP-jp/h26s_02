package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/traP-jp/h26s_02/domain"
)

type Post interface {
	CreatePost(ctx context.Context, id uuid.UUID, userName string) error
	GetPosts(ctx context.Context, referenceID uuid.UUID, referenceTime time.Time, limit int) ([]*domain.Post, error)
	GetPostByID(ctx context.Context, id uuid.UUID) (*domain.Post, error)
	GetPostsByUser(ctx context.Context, userName string) ([]*domain.Post, error)
	GetPostsByTags(ctx context.Context, tagNames []string) ([]*domain.Post, error)
}
