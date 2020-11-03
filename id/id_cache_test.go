package id

import (
	"testing"

	"github.com/hoaitan/cache/local"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	name := "abc"

	cache := New(local.New(local.Config{
		Enable: true,
		Size:   1000000,
	}), "test-id-cache").SetLoadFn(func(name string) (id string, err error) {
		return name, nil
	})
	assert.NotNil(t, cache)

	// Check empty cache
	ok, err := cache.IsExist(name)
	assert.Nil(t, err)
	assert.False(t, ok)

	// Load empty cache
	id, err := cache.GetOrSet(name)
	assert.Nil(t, err)
	assert.Equal(t, name, id)

	// Check loaded cache
	ok, err = cache.IsExist(name)
	assert.Nil(t, err)
	assert.True(t, ok)

	// Delete cache
	err = cache.Delete(name)
	assert.Nil(t, err)

	// Check deleted cache
	ok, err = cache.IsExist(name)
	assert.Nil(t, err)
	assert.False(t, ok)
}

func TestNew_With_Missing_LoadFn(t *testing.T) {
	name := "abc"
	cache := New(local.New(local.Config{
		Enable: true,
		Size:   1000000,
	}), "test-id-cache")

	// Load empty cache
	id, err := cache.GetOrSet(name)
	assert.NotNil(t, err)
	assert.Equal(t, "", id)
}
