package db

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type Tag struct {
	db *DB
}

func NewTag(db *DB) *Tag {
	return &Tag{
		db: db,
	}
}

func (t *Tag) CreatePostTags(ctx context.Context, postID uuid.UUID, tags []string) error {
	if len(tags) == 0 {
		return nil
	}

	type PostTag struct {
		PostID uuid.UUID `db:"post_id"`
		Tag    string    `db:"tag"`
	}
	postTags := make([]PostTag, 0, len(tags))
	for _, tag := range tags {
		postTags = append(postTags, PostTag{
			PostID: postID,
			Tag:    tag,
		})
	}

	_, err := t.db.db.NamedExecContext(ctx, "INSERT INTO `post_tags` (`post_id`, `name`) VALUES (:post_id, :tag)", postTags)
	if err != nil {
		return fmt.Errorf("create post tags: %w", err)
	}

	return nil
}
