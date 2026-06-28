package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/traP-jp/h26s_02/domain"
)

type Reaction interface {
	GetReactionCount(ctx context.Context, postID uuid.UUID) ([]*domain.ReactionCount, error)
	// post id の配列に紐づけられた reaction を取得する。
	GetReactionsByPostIDs(ctx context.Context, postIDs []uuid.UUID) (map[uuid.UUID][]*domain.ReactionCount, error)
}
