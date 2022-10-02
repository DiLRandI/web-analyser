package repository

import (
	"context"

	"github.com/DiLRandI/web-analyser/internal/dao"
)

type Results interface {
	Save(ctx context.Context, m *dao.Analyses) (int64, error)
	Remove(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (*dao.Analyses, error)
	GetAll(ctx context.Context) ([]*dao.Analyses, error)
}
