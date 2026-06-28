package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

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
		return fmt.Errorf("create post: %w", err)
	}

	return nil
}

func (p *Post) GetPosts(ctx context.Context, referenceTime time.Time, limit int) ([]*domain.Post, error) {
	if limit <= 0 {
		limit = 30
	} else if limit > 100 {
		limit = 100
	}

	query := `
		SELECT id, user_name, created_at
		FROM posts
		WHERE created_at < ?
		ORDER BY created_at DESC
		LIMIT ?`

	var records []posts
	if err := p.db.DB(ctx).SelectContext(ctx, &records, query, referenceTime, limit); err != nil {
		return nil, fmt.Errorf("get posts: %w", err)
	}

	posts := make([]*domain.Post, 0, len(records))
	for _, rec := range records {
		posts = append(posts, domain.NewPost(rec.ID, rec.UserName, rec.CreatedAt))
	}

	return posts, nil
}

func (p *Post) GetPostByID(ctx context.Context, id uuid.UUID) (*domain.Post, error) {
	var rec posts
	if err := p.db.DB(ctx).GetContext(ctx, &rec, "SELECT id, user_name, created_at FROM posts WHERE id = ?", id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrRecordNotFound
		}
		return nil, fmt.Errorf("get post by id: %w", err)
	}

	return domain.NewPost(rec.ID, rec.UserName, rec.CreatedAt), nil
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

	return domain.NewPost(post.ID, post.UserName, post.CreatedAt), nil
}

func (p *Post) GetPostsByUser(ctx context.Context, userName string) ([]*domain.Post, error) {
	var records []posts
	if err := p.db.DB(ctx).SelectContext(ctx, &records, "SELECT id, user_name, created_at FROM posts WHERE user_name = ? ORDER BY created_at DESC", userName); err != nil {
		return nil, fmt.Errorf("get posts by user: %w", err)
	}

	posts := make([]*domain.Post, 0, len(records))
	for _, rec := range records {
		posts = append(posts, domain.NewPost(rec.ID, rec.UserName, rec.CreatedAt))
	}

	return posts, nil
}
