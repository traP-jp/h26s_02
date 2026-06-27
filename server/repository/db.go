package repository

import "context"

type DB interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
