package lfu

import (
	"container/heap"
	"github.com/go-programming-tour-book/tourMemCache/cache"
	"github.com/go-programming-tour-book/tourMemCache/util"
	"log"
)

type lfu struct {
	// 缓存的最大容量 单位字节
	maxBytes int

	// 从缓存中移除节点时调用的回调函数
	onEvicted func(key string, value any)

	// 已使用字节数，不计算key
	usedBytes int

	queue *queue // 实现队列
	cache map[string]*entry
}

type entry struct {
	key   string
	value any
	freq  int
	index int
}

var _ cache.Cache = (*lfu)(nil)

func New(maxBytes int, onEvicted func(key string, value any)) cache.Cache {
	return &lfu{
		maxBytes:  maxBytes,
		onEvicted: onEvicted,
		usedBytes: 0,
		queue:     new(queue), // 可以提前分配较大内存减少 append 的拷贝
		cache:     make(map[string]*entry),
	}
}

func (c *lfu) delEntry(en *entry) {
	if en == nil {
		return
	}

	delete(c.cache, en.key)
	c.usedBytes -= util.CalcLen(en.value)
	if c.onEvicted != nil {
		c.onEvicted(en.key, en.value)
	}
}

func (c *lfu) Set(key string, value any) {
	if util.CalcLen(value) > c.maxBytes {
		log.Printf("超过缓存上界: maxBytes: %v, key: %v, value: %v", c.maxBytes, key, value)
		return
	}

	en, ok := c.cache[key]
	usedBytesGrow := util.CalcLen(value)
	if ok {
		usedBytesGrow -= util.CalcLen(en.value)
	}
	for c.usedBytes+usedBytesGrow > c.maxBytes {
		c.DelOldest()
	}

	if ok {
		c.usedBytes += usedBytesGrow
		c.queue.Update(en, value, en.freq+1)
		return
	}

	en = &entry{
		key:   key,
		value: value,
		freq:  1,
	}
	c.usedBytes += usedBytesGrow
	heap.Push(c.queue, en)
	c.cache[key] = en

}

func (c *lfu) Get(key string) any {
	en, ok := c.cache[key]
	if !ok {
		return nil
	}
	c.queue.Update(en, en.value, en.freq+1)
	return en.value
}

func (c *lfu) Del(key string) {
	en, ok := c.cache[key]
	if !ok {
		return
	}
	heap.Remove(c.queue, en.index)
	c.delEntry(en)
}

func (c *lfu) DelOldest() {
	if c.Len() == 0 {
		return
	}
	c.delEntry(heap.Pop(c.queue).(*entry))
}

func (c *lfu) Len() int {
	return c.queue.Len()
}
