package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/traP-jp/h26s_02/domain"
)

type Reaction interface {
	// すでに reaction がある場合は ErrUniqueKeyDuplicated。
	// 存在しない postID の場合は ErrViolatedForeignKey
	CreateReaction(ctx context.Context, postID uuid.UUID, reactionID int, userName string) error
	GetReactionCount(ctx context.Context, postID uuid.UUID) ([]*domain.ReactionCount, error)
}
