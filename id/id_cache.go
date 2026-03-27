package id

import (
	"context"
	"fmt"

	"github.com/hoaitan/cache"
)

const defaultTTL = 24 * 3600 // Long TTL: 1 day

type Cache interface {
	SetLoadFn(loadFn func(ctx context.Context, name string) (id string, err error)) Cache
	GetOrSet(ctx context.Context, name string) (id string, err error)
	IsExist(ctx context.Context, name string) (ok bool, err error)
	Delete(ctx context.Context, name string) (err error)
}

type idCache struct {
	cache     cache.Cache
	namespace string
	loadFn    func(ctx context.Context, name string) (id string, err error)
}

func New(cache cache.Cache, namespace string) Cache {
	return &idCache{
		cache:     cache,
		namespace: namespace,
	}
}

func (c *idCache) SetLoadFn(loadFn func(ctx context.Context, name string) (id string, err error)) Cache {
	c.loadFn = loadFn
	return c
}

func (c *idCache) GetOrSet(ctx context.Context, name string) (id string, err error) {
	err = c.cache.Get(ctx, cache.MakeKey(c.namespace, name), &id, func(ctx context.Context) error {
		if c.loadFn == nil {
			return fmt.Errorf("missing loadFn")
		}

		// Load ID
		id, err = c.loadFn(ctx, name)
		if err != nil {
			return err
		}

		// Set cache
		return c.cache.Set(ctx, cache.MakeKey(c.namespace, name), id, defaultTTL)
	})

	return id, err
}

func (c *idCache) IsExist(ctx context.Context, name string) (ok bool, err error) {
	return c.cache.IsExist(ctx, cache.MakeKey(c.namespace, name))
}

func (c *idCache) Delete(ctx context.Context, name string) (err error) {
	_, err = c.cache.Delete(ctx, cache.MakeKey(c.namespace, name))
	return err
}
