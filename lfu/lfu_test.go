package lfu

import (
	"github.com/matryer/is"
	"testing"
)

func TestSet(t *testing.T) {
	iss := is.New(t)

	cache := New(24, nil)
	cache.DelOldest()
	cache.Set("k1", 1)
	v := cache.Get("k1")
	iss.Equal(v, 1)

	cache.Del("k1")
	iss.Equal(0, cache.Len())
}

func TestOnEvicted(t *testing.T) {
	iss := is.New(t)

	keys := make([]string, 0, 8)
	onEvicted := func(key string, value any) {
		keys = append(keys, key)
	}
	cache := New(32, onEvicted)

	cache.Set("k1", 1)
	cache.Set("k2", 2)
	cache.Get("k1")
	cache.Get("k1")
	cache.Get("k2")
	cache.Set("k3", 3)
	cache.Set("k4", 4)

	expected := []string{"k2", "k3"}

	iss.Equal(expected, keys)
	iss.Equal(2, cache.Len())
}
