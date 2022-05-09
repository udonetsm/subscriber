package cache

import (
	"sync"
)

type Cache struct {
	sync.RWMutex
	data map[string]string
}

func New() *Cache {
	items := make(map[string]string)
	return &Cache{
		data: items,
	}
}

func (c *Cache) Set(key, value string) {
	c.Lock()
	defer c.Unlock()
	c.data[key] = value
}

func (c *Cache) Get(key string) (string, bool) {
	c.RLock()
	defer c.RUnlock()
	item, ok := c.data[key]
	if !ok {
		return "", false
	}
	return item, true
}
