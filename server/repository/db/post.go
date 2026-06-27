package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/traP-jp/h26s_02/domain"
	"github.com/traP-jp/h26s_02/repository"
)

type Post struct {
	db *DB
}

func NewPost(db *DB) *Post {
	return &Post{
		db: db,
	}
}

func (p *Post) CreatePost(ctx context.Context, id uuid.UUID, userName string) error {
	_, err := p.db.DB(ctx).
		ExecContext(ctx, "INSERT INTO `posts` (`id`, `user_name`) VALUES (?, ?)", id, userName)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}
func (p *Post) GetPost(ctx context.Context, id uuid.UUID) (*domain.Post, error) {
	var post posts
	err := p.db.DB(ctx).GetContext(ctx, &post, "SELECT  id, user_name, created_at FROM posts WHERE id = ?", id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, repository.ErrUniqueKeyDuplicated
	}
	if err != nil {
		return nil, fmt.Errorf("get post: %w", err)
	}

	return domain.NewPost(post.UserName, post.ID, post.CreatedAt), nil
}
