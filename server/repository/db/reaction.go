package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
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

func (r *Reaction) DeleteReaction(ctx context.Context, postID uuid.UUID, userName string, reactionID int) error {
	result, err := r.db.DB(ctx).ExecContext(ctx,
		"DELETE FROM post_reactions WHERE post_id = ? AND user_name = ? AND reaction_id = ?",
		postID,
		userName,
		reactionID,
	)
	if err != nil {
		return fmt.Errorf("delete reaction: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete reaction: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("no record deleted")
	}

	return nil
}
func (r *Reaction) GetReactionsByPostIDs(ctx context.Context, postIDs []uuid.UUID) (map[uuid.UUID][]*domain.ReactionCount, error) {
	if len(postIDs) == 0 {
		return make(map[uuid.UUID][]*domain.ReactionCount), nil
	}

	query := `SELECT post_id, reaction_id, COUNT(post_id) AS count FROM post_reactions WHERE post_id IN (?) GROUP BY post_id, reaction_id`
	query, args, err := sqlx.In(query, postIDs)
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var rows []struct {
		PostID     uuid.UUID `db:"post_id"`
		ReactionID int       `db:"reaction_id"`
		Count      int       `db:"count"`
	}

	if err := r.db.DB(ctx).SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, fmt.Errorf("select reactions: %w", err)
	}

	result := make(map[uuid.UUID][]*domain.ReactionCount)
	for _, row := range rows {
		// domain.NewReaction がある前提です。なければ構造体を直接初期化してください
		reaction := domain.NewReaction(row.ReactionID, row.Count)
		result[row.PostID] = append(result[row.PostID], reaction)
	}
	return result, nil
}
