package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/traP-jp/h26s_02/domain"
)

type Reaction interface {
	GetReactionCount(ctx context.Context, postID uuid.UUID) ([]*domain.ReactionCount, error)
}
