package db

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/traP-jp/h26s_02/domain"
)

type Reaction struct {
	db *DB
}

func NewReaction(db *DB) *Reaction {
	return &Reaction{
		db: db,
	}
}

func (r *Reaction) GetReactionCount(ctx context.Context, postID uuid.UUID) ([]*domain.ReactionCount, error) {
	type reactionCount struct {
		ReactionID int `db:"reaction_id"`
		Count      int `db:"count"`
	}
	var reactionCounts []reactionCount
	err := r.db.DB(ctx).SelectContext(ctx, &reactionCounts,
		"SELECT reaction_id, COUNT(*) AS count FROM post_reactions WHERE post_id = ? GROUP BY reaction_id",
		postID,
	)
	if err != nil {
		return nil, fmt.Errorf("get reactions: %w", err)
	}

	reactions := make([]*domain.ReactionCount, 0, len(reactionCounts))
	for _, reactionCount := range reactionCounts {
		reactions = append(reactions, domain.NewReaction(
			reactionCount.ReactionID,
			reactionCount.Count,
		))
	}

	return reactions, nil
}
