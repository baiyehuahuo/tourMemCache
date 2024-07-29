package fifo

import "testing"
import "github.com/matryer/is"

func TestSetGet(t *testing.T) {
	iss := is.New(t)

	cache := New(24, nil)
	cache.DelOldest()
	cache.Set("k1", 1)
	v := cache.Get("k1")
	iss.Equal(v, 1)
	cache.Del("k1")
	iss.Equal(cache.Get("k1"), nil)
	iss.Equal(cache.Len(), 0)
}

func TestOnEvicted(t *testing.T) {
	iss := is.New(t)

	keys := make([]string, 0, 8)

	cache := New(48, func(key string, value any) {
		keys = append(keys, key)
	})
	setKeys := []string{"k1", "k2", "k3", "k4", "k5"}
	cache.Set(setKeys[0], 1)
	cache.Set(setKeys[1], 1)
	cache.Set(setKeys[2], 1)
	cache.Set(setKeys[3], 1)
	cache.Set(setKeys[4], 1)
	iss.Equal(cache.Len(), 3)
	iss.Equal(keys, setKeys[:2])
	cache.DelOldest()
	cache.DelOldest()
	iss.Equal(keys, setKeys[:4])
}
