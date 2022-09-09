package lru

import (
	"container/list"
	"fmt"
	"strings"
	"sync"
)

type Item struct {
	KeyPtr *list.Element
	Data   any
}

type CacheLRU struct {
	Queue *list.List
	Cache map[any]*Item // key would be the type used in correspondent method

	mu sync.RWMutex

	capacity int
}

func NewCacheLRU(capacity int) *CacheLRU {
	return &CacheLRU{
		Queue:    list.New(),
		Cache:    make(map[any]*Item),
		capacity: capacity,
	}
}

func (c *CacheLRU) String() string {
	var res []string

	res = append(res, fmt.Sprintf("Capacity: %d", c.capacity))
	res = append(res, "Cached:")

	for key, item := range c.Cache {
		res = append(res, fmt.Sprintf("key: %v, value: %v", key, item.Data))
	}

	return strings.Join(res, "\n")
}

func (c *CacheLRU) Put(key, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, exists := c.Cache[key]; !exists {
		if c.capacity == len(c.Cache) {
			lru := c.Queue.Back()
			c.Queue.Remove(lru)

			delete(c.Cache, lru.Value.(int))
		}

		c.Cache[key] = &Item{
			KeyPtr: c.Queue.PushFront(key),
			Data:   value,
		}
	} else {
		node.Data = value

		c.Cache[key] = node
		c.Queue.MoveToFront(node.KeyPtr)
	}
}

func (c *CacheLRU) Get(key any) any {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if item, exists := c.Cache[key]; exists {
		c.Queue.MoveToFront(item.KeyPtr)

		return item.Data
	}

	return nil
}

func (c *CacheLRU) Delete(key any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, exists := c.Cache[key]
	if !exists {
		return
	}

	c.Queue.Remove(item.KeyPtr)
	delete(c.Cache, key)
}
