package mem

import (
	"context"
	"testing"

	"github.com/DiLRandI/web-analyser/internal/dao"
	"github.com/stretchr/testify/assert"
)

func Test_save_should_add_item_to_map(t *testing.T) {
	t.Cleanup(cleanup)
	sut := NewResultInMemory()
	assert.Empty(t, data)

	id, err := sut.Save(context.Background(), &dao.Analyses{})

	assert.NoError(t, err)
	assert.Greater(t, id, int64(0))
	assert.NotEmpty(t, data)
	assert.Len(t, data, 1)

	_, ok := data[id]
	assert.True(t, ok)
}

func Test_update_should_update_the_values_in_map(t *testing.T) {
	t.Cleanup(cleanup)
	data[1] = &dao.Analyses{
		Id:    1,
		Title: "test",
		Url:   "https://www.test.com",
	}
	currentId = 1
	sut := NewResultInMemory()
	assert.NotEmpty(t, data)

	updateModel := &dao.Analyses{
		Id:    1,
		Title: "test updated",
		Url:   "https://www.test-updated.com",
	}
	err := sut.Update(context.Background(), 1, updateModel)

	assert.NoError(t, err)
	assert.Equal(t, data[1], updateModel)
}

func Test_update_should_throw_an_error_for_invalid_id(t *testing.T) {
	t.Cleanup(cleanup)
	data[1] = &dao.Analyses{}
	currentId = 1

	sut := NewResultInMemory()
	invalidId := int64(2)
	err := sut.Update(context.Background(), invalidId, &dao.Analyses{})

	assert.ErrorIs(t, err, ResultNotFoundErr)
}

func Test_save_should_persis_two_consecutive_items(t *testing.T) {
	t.Cleanup(cleanup)
	sut := NewResultInMemory()
	assert.Empty(t, data)

	id1, err1 := sut.Save(context.Background(), &dao.Analyses{})

	assert.NoError(t, err1)
	assert.Greater(t, id1, int64(0))
	assert.NotEmpty(t, data)
	assert.Len(t, data, 1)

	_, ok1 := data[id1]
	assert.True(t, ok1)

	id2, err2 := sut.Save(context.Background(), &dao.Analyses{})

	assert.NoError(t, err2)
	assert.Greater(t, id2, int64(1))
	assert.NotEmpty(t, data)
	assert.Len(t, data, 2)

	_, ok2 := data[id2]
	assert.True(t, ok2)

}

func Test_remove_should_remove_the_item_from_map(t *testing.T) {
	t.Cleanup(cleanup)
	data[1] = &dao.Analyses{}
	currentId = 1
	sut := NewResultInMemory()
	assert.NotEmpty(t, data)

	err := sut.Remove(context.Background(), 1)

	assert.NoError(t, err)
	assert.Empty(t, data)
}

func Test_remove_should_remove_throw_an_error_for_invalid_id(t *testing.T) {
	t.Cleanup(cleanup)
	data[1] = &dao.Analyses{}
	currentId = 1

	sut := NewResultInMemory()
	invalidId := int64(2)
	err := sut.Remove(context.Background(), invalidId)

	assert.ErrorIs(t, err, ResultNotFoundErr)
}

func Test_get_should_return_item_on_map(t *testing.T) {
	t.Cleanup(cleanup)
	s := &dao.Analyses{}
	data[1] = s
	currentId = 1

	sut := NewResultInMemory()
	result, err := sut.Get(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, result, s)
}

func Test_get_should_throw_an_error_for_invalid_id(t *testing.T) {
	t.Cleanup(cleanup)
	s := &dao.Analyses{}
	data[1] = s
	currentId = 1

	invalidId := int64(2)
	sut := NewResultInMemory()
	result, err := sut.Get(context.Background(), invalidId)

	assert.ErrorIs(t, err, ResultNotFoundErr)
	assert.Nil(t, result)
}

func Test_get_all_return_all_the_items_in_map(t *testing.T) {
	t.Cleanup(cleanup)
	data[1] = &dao.Analyses{}
	data[2] = &dao.Analyses{}
	data[3] = &dao.Analyses{}
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
