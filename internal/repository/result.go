package repository

import "context"

type Results interface {
	Save(ctx context.Context, m any) (int64, error)
	Remove(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (any, error)
	GetAll(ctx context.Context) ([]any, error)
}
