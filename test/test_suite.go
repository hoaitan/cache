package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hoaitan/cache"
	"github.com/stretchr/testify/assert"
)

var (
	missCacheErr = fmt.Errorf("missing cache hit")
	missCacheFn  = func(ctx context.Context) error {
		return missCacheErr
	}
	nilFn = func(ctx context.Context) error {
		return nil
	}
)

func GetTestSuite(isEnableCache bool) []func(t *testing.T, c cache.Cache) {
	// Test suite for enable cache
	if isEnableCache {
		return []func(t *testing.T, c cache.Cache){
			testSet,
			testGet,
			testDelete,
			testIsExist,
			testFlush,
		}
	}

	// Test suite for disable cache
	return []func(t *testing.T, c cache.Cache){
		testDisableCacheSet,
		testDisableCacheGet,
		testDisableCacheDelete,
		testDisableCacheIsExist,
		testDisableCacheFlush,
	}
}

func testSet(t *testing.T, c cache.Cache) {
	ctx := context.Background()

	// Set empty key
	err := c.Set(ctx, "", "", 0)
	assert.Nil(t, err)

	// Set key without TTL
	err = c.Set(ctx, "test:set:ttl=0", 1, 0)
	assert.Nil(t, err)

	cacheInt := 0
	err = c.Get(ctx, "test:set:ttl=0", &cacheInt, missCacheFn)
	assert.Nil(t, err)
	assert.Equal(t, 1, cacheInt)

	// Set key with TTL
	err = c.Set(ctx, "test:set:ttl=10", 1, 1)
	assert.Nil(t, err)

	cacheInt = 0
	err = c.Get(ctx, "test:set:ttl=10", &cacheInt, missCacheFn)
	assert.Nil(t, err)
	assert.Equal(t, 1, cacheInt)

	// Set key with TTL and expired
	err = c.Set(ctx, "test:set:ttl=1", 1, 1)
	assert.Nil(t, err)

	time.Sleep(time.Second * 2)
	cacheInt = 0
	err = c.Get(ctx, "test:set:ttl=1", &cacheInt, missCacheFn)
	assert.Error(t, err)
}

func testGet(t *testing.T, c cache.Cache) {
	ctx := context.Background()

	// Trigger missing cache
	err := c.Get(ctx, "test:get:missing", nil, missCacheFn)
	assert.Equal(t, missCacheErr, err)

	// Skip missCacheFn if nil
	err = c.Get(ctx, "test:get:missing", nil, nilFn)
	assert.Nil(t, err)

	// Get normal cache
	c.Set(ctx, "test:set", 1, 0)

	cacheInt := 0
	err = c.Get(ctx, "test:set", &cacheInt, missCacheFn)
	assert.Nil(t, err)
	assert.Equal(t, 1, cacheInt)

	// Invalid type between set and get cache
	c.Set(ctx, "test:set:invalid", 1, 0)

	cacheString := ""
	err = c.Get(ctx, "test:set:invalid", &cacheString, missCacheFn)
	assert.Error(t, err)
	assert.NotEqual(t, missCacheErr, err)
}

func testDelete(t *testing.T, c cache.Cache) {
	ctx := context.Background()

	// Not found key
	ok, err := c.Delete(ctx, "test:delete:not-found")
	assert.False(t, ok)
	assert.Nil(t, err)

	// Exist key
	c.Set(ctx, "test:delete", 1, 0)
	ok, err = c.Delete(ctx, "test:delete")
	assert.True(t, ok)
	assert.Nil(t, err)
}

func testIsExist(t *testing.T, c cache.Cache) {
	ctx := context.Background()

	// Not found key
	ok, err := c.IsExist(ctx, "test:is-exist:not-found")
	assert.False(t, ok)
	assert.Nil(t, err)

	// Exist key
	c.Set(ctx, "test:is-exist", 1, 0)
	ok, err = c.IsExist(ctx, "test:is-exist")
	assert.True(t, ok)
	assert.Nil(t, err)

	// Expired key
	c.Set(ctx, "test:is-exist:ttl=1", 1, 1)

	time.Sleep(time.Second * 2)
	ok, err = c.IsExist(ctx, "test:is-exist:ttl=1")
	assert.False(t, ok)
	assert.Nil(t, err)
}
func testFlush(t *testing.T, c cache.Cache) {
	ctx := context.Background()

	c.Set(ctx, "test:flush", 1, 0)
	_, err := c.Flush(ctx)

	// Flush function can be not supported in Redis implementation
	if err == cache.NotSupportedErr {
		fmt.Println("Flush is not supported. Skip!!!")
		return
	}
	assert.Nil(t, err)

	ok, err := c.IsExist(ctx, "test:flush")
	assert.False(t, ok)
	assert.Nil(t, err)
}

// Disable cache tests
func testDisableCacheSet(t *testing.T, c cache.Cache) {
	// Set empty key
	err := c.Set(context.Background(), "test:set:disable", "123", 0)
	assert.Nil(t, err)
}

func testDisableCacheGet(t *testing.T, c cache.Cache) {
	// missCacheFn should be called and return missCacheErr
	err := c.Get(context.Background(), "test:get:disable", nil, missCacheFn)
	assert.Equal(t, missCacheErr, err)
}

func testDisableCacheDelete(t *testing.T, c cache.Cache) {
	ok, err := c.Delete(context.Background(), "test:delete:disable")
	assert.False(t, ok)
	assert.Nil(t, err)
}

func testDisableCacheIsExist(t *testing.T, c cache.Cache) {
	// Not found key
	ok, err := c.IsExist(context.Background(), "test:is-exist:disable")
	assert.False(t, ok)
	assert.Nil(t, err)
}
func testDisableCacheFlush(t *testing.T, c cache.Cache) {
	count, err := c.Flush(context.Background())

	// Flush function can be not supported in Redis implementation
	if err == cache.NotSupportedErr {
		fmt.Println("Flush is not supported. Skip!!!")
		return
	}
	assert.Nil(t, err)
	assert.Equal(t, 0, count)
}
