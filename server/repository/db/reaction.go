package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/traP-jp/h26s_02/domain"
	"github.com/traP-jp/h26s_02/repository"
)

type Reaction struct {
	db *DB
}

func NewReaction(db *DB) *Reaction {
	return &Reaction{
		db: db,
	}
}

func (r *Reaction) CreateReaction(ctx context.Context, postID uuid.UUID, reactionID int, userName string) error {
	_, err := r.db.DB(ctx).ExecContext(ctx,
		"INSERT INTO post_reactions (post_id, reaction_id, user_name, created_at) VALUES (?, ?, ?, NOW())",
		postID,
		reactionID,
		userName,
	)
	if err != nil {
		if mysqlErr, ok := errors.AsType[*mysql.MySQLError](err); ok {
			switch mysqlErr.Number {
			case 1062:
				return repository.ErrUniqueKeyDuplicated
			case 1452:
				return repository.ErrViolatedForeignKey
			}
		}
		return fmt.Errorf("create reaction: %w", err)
	}
	return nil
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
