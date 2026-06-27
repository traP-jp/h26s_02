package db

import (
	"time"

	"github.com/google/uuid"
)

type Posts struct {
	ID        uuid.UUID `db:"id"`
	UserName  string    `db:"user_name"`
	CreatedAt time.Time `db:"created_at"`
}

type PostTags struct {
	PostID uuid.UUID `db:"post_id"`
	Name   string    `db:"name"`
}

type PostReactions struct {
	PostID     uuid.UUID `db:"post_id"`
	ReactionID int       `db:"reaction_id"`
	UserName   string    `db:"user_name"`
	CreatedAt  time.Time `db:"created_at"`
}
