package repository

import (
	"context"

	"github.com/google/uuid"
)

type Post interface {
	CreatePost(ctx context.Context, id uuid.UUID, userName string) error
}
