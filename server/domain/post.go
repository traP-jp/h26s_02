package domain

import (
	"github.com/google/uuid"
	"time"
)

type Post struct {
	id        string
	userName  string
	createdAt string
}

func NewPost(userName string) *Post {
	return &Post{
		id:        uuid.New().String(),
		userName:  userName,
		createdAt: time.Now().Format(time.RFC3339),
	}
}
func (p *Post) GetID() string {
	return p.id
}
func (p *Post) GetUserName() string {
	return p.userName
}
func (p *Post) GetCreatedAt() string {
	return p.createdAt
}
