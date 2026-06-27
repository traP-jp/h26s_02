package db

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/traP-jp/h26s_02/domain"
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

type postRecord struct {
	ID        uuid.UUID `db:"id"`
	UserName  string    `db:"user_name"`
	CreatedAt time.Time `db:"created_at"`
}

func (p *Post) GetPosts(referenceTime time.Time, limit int) ([]*domain.Post, error) {
	if limit <= 0 {
		limit = 30
	}

	query := `
		SELECT id, user_name, created_at
		FROM posts
		WHERE created_at < ?
		ORDER BY created_at DESC
		LIMIT ?`

	var records []postRecord
	if err := p.db.db.Select(&records, query, referenceTime, limit); err != nil {
		return nil, fmt.Errorf("get posts: %w", err)
	}

	posts := make([]*domain.Post, 0, len(records))
	for _, rec := range records {
		posts = append(posts, domain.NewPost(rec.ID, rec.UserName, rec.CreatedAt))
	}

	return posts, nil
}

func (p *Post) GetPostByID(id string) (*domain.Post, error) {
	postID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("parse post id: %w", err)
	}

	var rec postRecord
	if err := p.db.db.Get(&rec, "SELECT id, user_name, created_at FROM posts WHERE id = ?", postID); err != nil {
		return nil, fmt.Errorf("get post by id: %w", err)
	}

	return domain.NewPost(rec.ID, rec.UserName, rec.CreatedAt), nil
}
