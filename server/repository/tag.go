package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/traP-jp/h26s_02/domain"
)

type Tag interface {
	CreatePostTags(ctx context.Context, postID uuid.UUID, tags []string) error

	// タグ一覧の取得
	GetTags(ctx context.Context) ([]TagCount, error)
	GetPostTags(ctx context.Context, postID uuid.UUID) ([]string, error)
}

type TagCount struct {
	Tag   domain.Tag
	Count int
}
