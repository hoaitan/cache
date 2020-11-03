package multi

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"

	"github.com/hoaitan/cache"
	"github.com/hoaitan/cache/local"
	"github.com/hoaitan/cache/redis"
	"github.com/hoaitan/cache/test"
)

func TestCacheImplement(t *testing.T) {
	caches := []cache.Cache{
		// Disable local cache
		local.New(local.Config{
			Enable: false,
			Size:   1000000,
		}),

		// Disable Redis cache
		redis.New(redis.Config{
			Enable:   false,
			Endpoint: "localhost:6379",
		}, "test"),

		// Enable local cache
		local.New(local.Config{
			Enable: true,
			Size:   1000000,
		}),
	}

	enableRedisCache := redis.New(redis.Config{
		Enable:   true,
		Endpoint: "localhost:6379",
	}, "test")

	// Is Redis ready for testing
	if enableRedisCache.IsReady() {
		caches = append(caches, enableRedisCache)
	}

	c := New(caches...)

	for _, fn := range test.GetTestSuite(true) {
		t.Run(fmt.Sprintf("fn=%s", runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()), func(t *testing.T) {
			fn(t, c)
		})

	}

}
