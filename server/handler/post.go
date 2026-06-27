package handler

import "github.com/traP-jp/h26s_02/repository"

type Post struct {
	postRepository repository.Post
}

func NewPost(postRepository repository.Post) *Post {
	return &Post{
		postRepository: postRepository,
	}
}
