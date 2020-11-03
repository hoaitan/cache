package redis

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"

	"github.com/hoaitan/cache/test"
)

func TestCacheImplement_Enable(t *testing.T) {
	c := New(Config{
		Enable:     true,
		Endpoint:   "localhost:6379",
		Timeout:    60,
		DefaultTTL: 60,
	}, "test")

	// Is Redis ready for testing
	if !c.IsReady() {
		fmt.Println("Redis is not ready for testing. Exist!!!")
		return
	}

	for _, fn := range test.GetTestSuite(true) {
		t.Run(fmt.Sprintf("fn=%s", runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()), func(t *testing.T) {
			fn(t, c)
		})
	}
}

func TestCacheImplement_Disable(t *testing.T) {
	c := New(Config{
		Enable:     false,
		Endpoint:   "localhost:6379",
		Timeout:    60,
		DefaultTTL: 60,
	}, "test")

	// Is Redis ready for testing
	if !c.IsReady() {
		fmt.Println("Redis is not ready for testing. Exist!!!")
		return
	}

	for _, fn := range test.GetTestSuite(false) {
		t.Run(fmt.Sprintf("fn=%s", runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()), func(t *testing.T) {
			fn(t, c)
		})
	}
}
