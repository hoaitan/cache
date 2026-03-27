package multi

import (
	"context"
	"fmt"
	"strings"

	"github.com/hoaitan/cache"
)

type multiCaches struct {
	caches []cache.Cache
}

// Multi caches support cache in multi cache implements, order is important
func New(caches ...cache.Cache) cache.Cache {
	return &multiCaches{
		caches: caches,
	}
}

// Set caches for all implements
func (c *multiCaches) Set(ctx context.Context, key string, data interface{}, ttl int) (err error) {
	for _, cache := range c.caches {
		if err = cache.Set(ctx, key, data, ttl); err != nil {
			return err
		}
	}

	return nil
}

// Get first found cache in all implements
func (c *multiCaches) Get(ctx context.Context, key string, ptr interface{}, fn cache.MissCacheFn) (err error) {
	for _, cache := range c.caches {
		isFound := true
		checkFn := func(ctx context.Context) error {
			isFound = false
			return nil
		}

		if err = cache.Get(ctx, key, ptr, checkFn); err != nil {
			return err
		}

		if isFound {
			return nil
		}
	}

	// Call missing fn
	if fn == nil {
		return nil
	}

	return fn(ctx)
}

// Delete cache in all implements
func (c *multiCaches) Delete(ctx context.Context, key string) (ok bool, err error) {
	for _, cache := range c.caches {
		_ok, err := cache.Delete(ctx, key)
		ok = ok || _ok

		if err != nil {
			return ok, err
		}
	}

	return ok, nil
}

// Get first found cache in all implements
func (c *multiCaches) IsExist(ctx context.Context, key string) (ok bool, err error) {
	for _, cache := range c.caches {
		if ok, err = cache.IsExist(ctx, key); ok || err != nil {
			return ok, err
		}
	}

	return false, nil
}

// Flush all implements
func (c *multiCaches) Flush(ctx context.Context) (count int, err error) {
	for _, cache := range c.caches {
		_count, err := cache.Flush(ctx)
		if err != nil {
			return count, err
		}

		count = +_count
	}

	return count, nil
}

// Check all cache implements are ready or not
func (c *multiCaches) IsReady(ctx context.Context) (ok bool) {
	for _, cache := range c.caches {
		if ok = cache.IsReady(ctx); !ok {
			return false
		}
	}

	return true
}

// IsEnable is true if there is an enable cache
func (c *multiCaches) IsEnable() (ok bool) {
	for _, cache := range c.caches {
		if ok = cache.IsEnable(); ok {
			return true
		}
	}

	return false
}

// Close all cache implements
func (c *multiCaches) Close() error {
	var errS []string
	for _, cache := range c.caches {
		if err := cache.Close(); err != nil {
			errS = append(errS, err.Error())
		}
	}

	if len(errS) == 0 {
		return nil
	}

	return fmt.Errorf("errors: %s", strings.Join(errS, ", "))
}
