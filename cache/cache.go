package cache

type Cache interface {
	Set(key string, value any)
	Get(key string) any
	Del(key string)
	DelOldest()
	Len() int
}

type Value interface {
	Len() int
}
