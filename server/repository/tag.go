package repository

import (
	"context"

	"github.com/google/uuid"
)

type Tag interface {
	CreatePostTags(ctx context.Context, postID uuid.UUID, tags []string) error
}
