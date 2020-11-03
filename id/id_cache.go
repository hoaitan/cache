package id

import (
	"fmt"

	"github.com/hoaitan/cache"
)

const defaultTTL = 24 * 3600 // Long TTL: 1 day

type Cache interface {
	SetLoadFn(loadFn func(name string) (id string, err error)) Cache
	GetOrSet(name string) (id string, err error)
	IsExist(name string) (ok bool, err error)
	Delete(name string) (err error)
}

type idCache struct {
	cache     cache.Cache
	namespace string
	loadFn    func(name string) (id string, err error)
}

func New(cache cache.Cache, namespace string) Cache {
	return &idCache{
		cache:     cache,
		namespace: namespace,
	}
}

func (c *idCache) SetLoadFn(loadFn func(name string) (id string, err error)) Cache {
	c.loadFn = loadFn
	return c
}

func (c *idCache) GetOrSet(name string) (id string, err error) {
	err = c.cache.Get(cache.MakeKey(c.namespace, name), &id, func() error {
		if c.loadFn == nil {
			return fmt.Errorf("missing loadFn")
		}

		// Load ID
		id, err = c.loadFn(name)
		if err != nil {
			return err
		}

		// Set cache
		return c.cache.Set(cache.MakeKey(c.namespace, name), id, defaultTTL)
	})

	return id, err
}

func (c *idCache) IsExist(name string) (ok bool, err error) {
	return c.cache.IsExist(cache.MakeKey(c.namespace, name))
}

func (c *idCache) Delete(name string) (err error) {
	_, err = c.cache.Delete(cache.MakeKey(c.namespace, name))
	return err
}
