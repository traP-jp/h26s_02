package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/traP-jp/h26s_02/domain"
)

type Post interface {
	CreatePost(ctx context.Context, id uuid.UUID, userName string) error
	GetPost(ctx context.Context, id uuid.UUID) (*domain.Post, error)
}
