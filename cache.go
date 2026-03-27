package cache

import (
	"context"
	"fmt"
	"strings"
)

type MissCacheFn func(ctx context.Context) error

var NotSupportedErr = fmt.Errorf("not suppported")

type Cache interface {
	// ttl=-1: will use default TTL
	// ttl=0 : no expire
	Set(ctx context.Context, key string, data interface{}, ttl int) (err error)
	Get(ctx context.Context, key string, ptr interface{}, fn MissCacheFn) (err error)
	Delete(ctx context.Context, key string) (ok bool, err error)
	IsExist(ctx context.Context, key string) (ok bool, err error)
	Flush(ctx context.Context) (count int, err error)
	IsReady(ctx context.Context) (ok bool)
	IsEnable() (ok bool)
	Close() (err error)
}

func MakeKey(parts ...string) string {
	return strings.Join(parts, ":")
}
