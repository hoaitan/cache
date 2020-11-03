package local

import (
	"bytes"
	"encoding/gob"

	"github.com/coocood/freecache"
	"github.com/hoaitan/cache"
)

type localCache struct {
	cacheEngine *freecache.Cache
	cf          Config
}

// New local cache with size (byte, min = 512KB)
func New(cf Config) cache.Cache {
	return &localCache{
		cacheEngine: freecache.NewCache(cf.Size),
		cf:          cf,
	}
}

func (c *localCache) Set(key string, data interface{}, ttl int) (err error) {
	if !c.IsEnable() {
		return nil
	}
	if ttl < 0 {
		ttl = c.cf.DefaultTTL
	}

	// Encode data
	b, err := encode(data)
	if err != nil {
		return err
	}

	// Set value to cache engine
	return c.cacheEngine.Set([]byte(key), b, ttl)
}

func (c *localCache) Get(key string, ptr interface{}, fn cache.MissCacheFn) (err error) {
	if !c.IsEnable() {
		return fn()
	}

	// Get cached value
	v, err := c.cacheEngine.Get([]byte(key))
	if err != nil {
		// Call function if missing cache
		if fn == nil {
			return nil
		}

		return fn()
	}

	// Decode
	if err = decode(v, ptr); err != nil {
		return err
	}

	return nil
}

func (c *localCache) Delete(key string) (ok bool, err error) {
	if !c.IsEnable() {
		return false, nil
	}

	return c.cacheEngine.Del([]byte(key)), nil

}

func (c *localCache) IsExist(key string) (ok bool, err error) {
	if !c.IsEnable() {
		return false, nil
	}

	if _, err := c.cacheEngine.Get([]byte(key)); err != nil {
		return false, nil
	}

	return true, nil
}

func (c *localCache) Flush() (count int, err error) {
	if !c.IsEnable() {
		return 0, nil
	}

	count = int(c.cacheEngine.EntryCount())
	c.cacheEngine.Clear()
	c.cacheEngine.ResetStatistics()

	return count, nil
}

func (c *localCache) IsReady() bool {
	return true
}

func (c *localCache) IsEnable() bool {
	return c.cf.Enable
}

func (c *localCache) Close() error {
	return nil
}

func encode(data interface{}) ([]byte, error) {
	buff := new(bytes.Buffer)
	enc := gob.NewEncoder(buff)

	err := enc.Encode(data)
	return buff.Bytes(), err
}

func decode(data []byte, ptr interface{}) error {
	buff := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buff)

	err := dec.Decode(ptr)
	return err
}
