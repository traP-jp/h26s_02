package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/traP-jp/h26s_02/domain"
)

type Post interface {
	CreatePost(ctx context.Context, id uuid.UUID, userName string) error
	GetPosts(referenceTime time.Time, limit int) ([]*domain.Post, error)
	GetPostByID(id string) (*domain.Post, error)
	GetPost(ctx context.Context, id uuid.UUID) (*domain.Post, error)
}
