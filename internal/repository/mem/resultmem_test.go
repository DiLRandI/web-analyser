package mem

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_save_should_add_item_to_map(t *testing.T) {
	t.Cleanup(cleanup)
	sut := NewResultInMemory()
	assert.Empty(t, data)

	id, err := sut.Save(context.Background(), "Test data")

	assert.NoError(t, err)
	assert.Greater(t, id, int64(0))
	assert.NotEmpty(t, data)
	assert.Len(t, data, 1)

	_, ok := data[id]
	assert.True(t, ok)
}

func Test_save_should_persis_two_consecutive_items(t *testing.T) {
	t.Cleanup(cleanup)
	sut := NewResultInMemory()
	assert.Empty(t, data)

	id1, err1 := sut.Save(context.Background(), "Test data")

	assert.NoError(t, err1)
	assert.Greater(t, id1, int64(0))
	assert.NotEmpty(t, data)
	assert.Len(t, data, 1)

	_, ok1 := data[id1]
	assert.True(t, ok1)

	id2, err2 := sut.Save(context.Background(), "Test data")

	assert.NoError(t, err2)
	assert.Greater(t, id2, int64(1))
	assert.NotEmpty(t, data)
	assert.Len(t, data, 2)

	_, ok2 := data[id2]
	assert.True(t, ok2)

}

func Test_remove_should_remove_the_item_from_map(t *testing.T) {
	t.Cleanup(cleanup)
	data[1] = "Test data"
	currentId = 1
	sut := NewResultInMemory()
	assert.NotEmpty(t, data)

	err := sut.Remove(context.Background(), 1)

	assert.NoError(t, err)
	assert.Empty(t, data)
}

func Test_remove_should_remove_throw_an_error_for_invalid_id(t *testing.T) {
	t.Cleanup(cleanup)
	data[1] = "Test data"
	currentId = 1

	sut := NewResultInMemory()
	invalidId := int64(2)
	err := sut.Remove(context.Background(), invalidId)

	assert.ErrorContains(t, err,
		fmt.Sprintf("Unable to remove the item with id %d, item not found", invalidId))
}

func Test_get_should_return_item_on_map(t *testing.T) {
	t.Cleanup(cleanup)
	s := "Test data"
	data[1] = s
	currentId = 1

	sut := NewResultInMemory()
	result, err := sut.Get(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Contains(t, result, s)
}

func Test_get_should_throw_an_error_for_invalid_id(t *testing.T) {
	t.Cleanup(cleanup)
	s := "Test data"
	data[1] = s
	currentId = 1

	invalidId := int64(2)
	sut := NewResultInMemory()
	result, err := sut.Get(context.Background(), invalidId)

	assert.ErrorContains(t, err,
		fmt.Sprintf("Unable to get the item with id %d, item not found", invalidId))
	assert.Nil(t, result)
}

func Test_get_all_return_all_the_items_in_map(t *testing.T) {
	t.Cleanup(cleanup)
	data[1] = "Test data 1"
	data[2] = "Test data 2"
	data[3] = "Test data 3"
	currentId = 3

	sut := NewResultInMemory()
	results, err := sut.GetAll(context.Background())

	assert.NoError(t, err)
	assert.Len(t, results, 3)
}

func cleanup() {
	currentId = 0
	for k := range data {
		delete(data, k)
	}

}
