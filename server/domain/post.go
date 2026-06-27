package domain

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	id        uuid.UUID
	userName  string
	tags      []string
	reactions []Reaction
	createdAt time.Time
}

func NewPost(id uuid.UUID, userName string, createdAt time.Time) *Post {
	return &Post{
		id:        id,
		userName:  userName,
		createdAt: createdAt,
	}
}
func (p *Post) GetID() uuid.UUID {
	return p.id
}
func (p *Post) GetUserName() string {
	return p.userName
}
func (p *Post) GetCreatedAt() time.Time {
	return p.createdAt
}
func (p *Post) GetReactions() []Reaction {
	return p.reactions
}
func (p *Post) GetTags() []string {
	return p.tags
}
