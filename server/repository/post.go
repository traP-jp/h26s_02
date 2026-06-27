package repository

import (
	"github.com/google/uuid"
	"github.com/traP-jp/h26s_02/domain"
)

type Post interface {
	// TODO
	GetPost(id uuid.UUID) (*domain.Post, error)
}
