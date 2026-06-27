package db

import (
	"time"

	"github.com/google/uuid"
)

type posts struct { // nolint:unused
	ID        uuid.UUID `db:"id"`
	UserName  string    `db:"user_name"`
	CreatedAt time.Time `db:"created_at"`
}

type postTags struct { // nolint:unused
	PostID uuid.UUID `db:"post_id"`
	Name   string    `db:"name"`
}

type postReactions struct { // nolint:unused
	PostID     uuid.UUID `db:"post_id"`
	ReactionID int       `db:"reaction_id"`
	UserName   string    `db:"user_name"`
	CreatedAt  time.Time `db:"created_at"`
}
