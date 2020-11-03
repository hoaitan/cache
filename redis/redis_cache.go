package redis

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	redisv8 "github.com/go-redis/redis/v8"
	"github.com/hoaitan/cache"
)

type redisCache struct {
	cacheEngine *redisv8.Client
	cf          Config
	keyPrefix   string
}

// New Redis cache
func New(cf Config, keyPrefix string) cache.Cache {
	return &redisCache{
		cacheEngine: redisv8.NewClient(&redisv8.Options{
			Addr:         cf.Endpoint,
			DialTimeout:  time.Duration(cf.Timeout) * time.Second,
			ReadTimeout:  time.Duration(cf.Timeout) * time.Second,
			WriteTimeout: time.Duration(cf.Timeout) * time.Second,
		}),
		cf:        cf,
		keyPrefix: strings.TrimRight(keyPrefix, ":"),
	}
}

func (c *redisCache) Set(key string, data interface{}, ttl int) (err error) {
	if !c.IsEnable() {
		return nil
	}
	if ttl < 0 {
		ttl = c.cf.DefaultTTL
	}

	// Encode data
	data, err = json.Marshal(data)
	if err != nil {
		return err
	}

	// Set value to cache engine
	return c.cacheEngine.Set(context.Background(), c.getKey(key), data, time.Duration(ttl)*time.Second).Err()
}

func (c *redisCache) Get(key string, ptr interface{}, fn cache.MissCacheFn) (err error) {
	if !c.IsEnable() {
		return fn()
	}

	// Get cached value
	v, err := c.cacheEngine.Get(context.Background(), c.getKey(key)).Result()
	if err != nil {
		// Call function if missing cache
		if fn == nil {
			return nil
		}
		return fn()
	}

	// Decode
	if err = json.Unmarshal([]byte(v), ptr); err != nil {
		return err
	}

	return nil
}

func (c *redisCache) Delete(key string) (ok bool, err error) {
	if !c.IsEnable() {
		return false, nil
	}

	count, err := c.cacheEngine.Del(context.Background(), c.getKey(key)).Result()

	return count > 0, err
}

func (c *redisCache) IsExist(key string) (ok bool, err error) {
	if !c.IsEnable() {
		return false, nil
	}

	count, err := c.cacheEngine.Exists(context.Background(), c.getKey(key)).Result()

	return count > 0, err
}

// Because performance issue, don't support this feature
func (c *redisCache) Flush() (count int, err error) {
	return 0, cache.NotSupportedErr
}

func (c *redisCache) IsReady() bool {
	if !c.IsEnable() {
		return true
	}

	if _, err := c.cacheEngine.Ping(context.Background()).Result(); err != nil {
		return false
	}

	return true
}

func (c *redisCache) IsEnable() bool {
	return c.cf.Enable
}

func (c *redisCache) Close() error {
	if !c.IsEnable() {
		return nil
	}

	return c.cacheEngine.Close()
}

func (c *redisCache) getKey(key string) string {
	if c.keyPrefix == "" {
		return key
	}

	return c.keyPrefix + ":" + key
}
