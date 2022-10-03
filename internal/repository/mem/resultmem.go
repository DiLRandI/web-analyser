package mem

import (
	"context"
	"sync/atomic"

	"github.com/DiLRandI/web-analyser/internal/dao"
	"github.com/DiLRandI/web-analyser/internal/repository"
)

var data map[int64]*dao.Analyses = make(map[int64]*dao.Analyses)
var currentId int64 = 0

type resultInMem struct {
}

func NewResultInMemory() repository.Results {
	return &resultInMem{}
}

func nextId() int64 {
	return atomic.AddInt64(&currentId, 1)
}

func (r *resultInMem) Save(ctx context.Context, m *dao.Analyses) (int64, error) {
	id := nextId()
	item := *m
	item.Id = id
	data[id] = &item
	return id, nil
}

func (r *resultInMem) Update(ctx context.Context, id int64, m *dao.Analyses) error {
	_, ok := data[id]
	if !ok {
		return ResultNotFoundErr
	}

	item := *m
	item.Id = id
	data[id] = &item
	return nil
}

func (r *resultInMem) Remove(ctx context.Context, id int64) error {
	if _, ok := data[id]; !ok {
		return ResultNotFoundErr
	}

	delete(data, id)
	return nil
}

func (r *resultInMem) Get(ctx context.Context, id int64) (*dao.Analyses, error) {
	if _, ok := data[id]; !ok {
		return nil, ResultNotFoundErr
	}
	item := *data[id]
	return &item, nil
}

func (r *resultInMem) GetAll(ctx context.Context) ([]*dao.Analyses, error) {
	results := []*dao.Analyses{}

	for _, d := range data {
		item := *d
		results = append(results, &item)
	}

	return results, nil
}
