package mem

import (
	"context"
	"fmt"
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
	data[id] = m
	return id, nil
}

func (r *resultInMem) Remove(ctx context.Context, id int64) error {
	if _, ok := data[id]; !ok {
		return fmt.Errorf("Unable to remove the item with id %d, item not found", id)
	}

	delete(data, id)
	return nil
}

func (r *resultInMem) Get(ctx context.Context, id int64) (*dao.Analyses, error) {
	if _, ok := data[id]; !ok {
		return nil, fmt.Errorf("Unable to get the item with id %d, item not found", id)
	}

	return data[id], nil
}

func (r *resultInMem) GetAll(ctx context.Context) ([]*dao.Analyses, error) {
	results := []*dao.Analyses{}

	for _, d := range data {
		results = append(results, d)
	}

	return results, nil
}
