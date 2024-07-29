package fifo

import (
	"container/list"
	"github.com/go-programming-tour-book/tourMemCache/cache"
	"github.com/go-programming-tour-book/tourMemCache/util"
	"log"
)

var _ cache.Cache = (*fifo)(nil)

// fifo 是一个 FIFO 缓存 但不是并发安全的
type fifo struct {
	// 缓存的最大容量 单位字节
	maxBytes int

	// 从缓存中移除节点时调用的回调函数
	onEvicted func(key string, value any)

	// 已使用字节数，不计算key
	usedBytes int

	ll    *list.List // 实现队列
	cache map[string]*list.Element
}

func New(maxBytes int, onEvicted func(key string, value any)) cache.Cache {
	return &fifo{
		maxBytes:  maxBytes,
		onEvicted: onEvicted,
		usedBytes: 0,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
	}
}

func (c *fifo) Set(key string, value any) {
	if util.CalcLen(value) > c.maxBytes {
		log.Printf("超过缓存上界: maxBytes: %v, key: %v, value: %v", c.maxBytes, key, value)
		return
	}

	// 缓存已存在的情况下
	if e, ok := c.cache[key]; ok {
		c.ll.MoveToBack(e)
		en := e.Value.(*cache.Entry)
		c.usedBytes += util.CalcLen(value) - util.CalcLen(en.Value)
		en.Value = value
	} else {
		// 缓存不存在的情况下
		entry := &cache.Entry{
			Key:   key,
			Value: value,
		}
		c.cache[key] = c.ll.PushBack(entry)
		c.usedBytes += util.CalcLen(entry.Value)
	}

	for c.usedBytes > c.maxBytes {
		c.DelOldest()
	}
}

func (c *fifo) Get(key string) any {
	e, ok := c.cache[key]
	if !ok {
		return nil
	}
	return e.Value.(*cache.Entry).Value
}

func (c *fifo) delElement(e *list.Element) {
	if e == nil {
		return
	}

	c.ll.Remove(e)
	delete(c.cache, e.Value.(*cache.Entry).Key)
	c.usedBytes -= util.CalcLen(e.Value.(*cache.Entry).Value)
	if c.onEvicted != nil {
		c.onEvicted(e.Value.(*cache.Entry).Key, e.Value.(*cache.Entry).Value)
	}
}

func (c *fifo) Del(key string) {
	e, ok := c.cache[key]
	if !ok {
		return
	}
	c.delElement(e)
}

func (c *fifo) DelOldest() {
	c.delElement(c.ll.Front())
}

func (c *fifo) Len() int {
	return len(c.cache)
}
