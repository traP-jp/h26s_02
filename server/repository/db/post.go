package db

import (
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
func (p *Post) GetPost(id uuid.UUID) (*domain.Post, error) {
	var post domain.Post
	err := p.db.db.Get(&post, "SELECT id, user_name, created_at FROM posts WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &post, nil
}
