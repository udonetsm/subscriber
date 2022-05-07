package cache

import (
	"sync"
)

type Cache struct {
	sync.RWMutex
	data map[string]Item
}

type Item struct {
	Val string
}

func New() *Cache {
	items := make(map[string]Item)
	return &Cache{
		data: items,
	}
}

func (c *Cache) Set(key, value string) {
	c.Lock()
	defer c.Unlock()
	c.data[key] = Item{
		Val: value,
	}
}

func (c *Cache) Get(key string) (string, bool) {
	c.RLock()
	defer c.RUnlock()
	item, ok := c.data[key]
	if !ok {
		return "", false
	}
	return item.Val, true
}
