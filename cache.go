package cache

import (
	"fmt"
	"strings"
)

type MissCacheFn func() error

var NotSupportedErr = fmt.Errorf("not suppported")

type Cache interface {
	// ttl=-1: will use default TTL
	// ttl=0 : no expire
	Set(key string, data interface{}, ttl int) (err error)
	Get(key string, ptr interface{}, fn MissCacheFn) (err error)
	Delete(key string) (ok bool, err error)
	IsExist(key string) (ok bool, err error)
	Flush() (count int, err error)
	IsReady() (ok bool)
	IsEnable() (ok bool)
	Close() (err error)
}

func MakeKey(parts ...string) string {
	return strings.Join(parts, ":")
}
