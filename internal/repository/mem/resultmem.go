package mem

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/DiLRandI/web-analyser/internal/repository"
)

var data map[int64]any = make(map[int64]any)
var currentId int64 = 0

type resultInMem struct {
	lock sync.Mutex
}

func NewResultInMemory() repository.Results {
	return &resultInMem{}
}

func nextId() int64 {
	return atomic.AddInt64(&currentId, 1)
}

func (r *resultInMem) Save(ctx context.Context, m any) (int64, error) {
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

func (r *resultInMem) Get(ctx context.Context, id int64) (any, error) {
	if _, ok := data[id]; !ok {
		return nil, fmt.Errorf("Unable to get the item with id %d, item not found", id)
	}

	return data[id], nil
}

func (r *resultInMem) GetAll(ctx context.Context) ([]any, error) {
	results := []any{}

	for _, d := range data {
		results = append(results, d)
	}

	return results, nil
}
