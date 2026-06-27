package domain

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	id        uuid.UUID
	userName  string
	createdAt time.Time
}

func NewPost(userName string, id uuid.UUID, createdAt time.Time) *Post {
	return &Post{
		id:        uuid.New(),
		userName:  userName,
		createdAt: time.Now(),
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
