package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeKey(t *testing.T) {
	key1 := MakeKey("this", "is", "key1")
	key1Copy := MakeKey("this", "is", "key1")
	key2 := MakeKey("this", "is", "key2")

	assert.NotEmpty(t, key1)
	assert.Equal(t, key1, key1Copy)
	assert.NotEqual(t, key1, key2)
}
