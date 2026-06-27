package db

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/traP-jp/h26s_02/domain"
	"github.com/traP-jp/h26s_02/repository"
)

type Tag struct {
	db *DB
}

func NewTag(db *DB) *Tag {
	return &Tag{
		db: db,
	}
}
func (t *Tag) GetPostTags(ctx context.Context, postID uuid.UUID) ([]string, error) {
	var postTags []postTags
	err := t.db.DB(ctx).SelectContext(ctx, &postTags, "SELECT name FROM post_tags WHERE post_id = ?", postID)
	if err != nil {
		return nil, fmt.Errorf("get post tags: %w", err)
	}
	tags := make([]string, 0, len(postTags))
	for _, postTag := range postTags {
		tags = append(tags, postTag.Name)
	}
	return tags, nil
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

	_, err := t.db.DB(ctx).NamedExecContext(ctx, "INSERT INTO `post_tags` (`post_id`, `name`) VALUES (:post_id, :tag)", postTags)
	if err != nil {
		return fmt.Errorf("create post tags: %w", err)
	}

	return nil
}

func (t *Tag) GetTags(ctx context.Context) ([]repository.TagCount, error) {
	type TagCount struct {
		Name  string `db:"name"`
		Count int    `db:"count"`
	}

	var tagCounts []TagCount
	err := t.db.DB(ctx).SelectContext(ctx, &tagCounts,
		"SELECT `name`, COUNT(`name`) as `count` FROM `post_tags` GROUP BY `name` ORDER BY `count` DESC")
	if err != nil {
		return nil, fmt.Errorf("get tags: %w", err)
	}

	tagStats := make([]repository.TagCount, 0, len(tagCounts))
	for _, tagCount := range tagCounts {
		tagStats = append(tagStats, repository.TagCount{
			Tag:   *domain.NewTag(tagCount.Name),
			Count: tagCount.Count,
		})
	}

	return tagStats, nil
}
