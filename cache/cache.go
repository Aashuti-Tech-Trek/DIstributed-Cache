package cache

import "sync"

type Cache struct {
	mu    sync.RWMutex
	store *LRUCache
}

func NewCache(capacity int) *Cache {
	return &Cache{
		store: NewLRU(capacity),
	}
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.store.Get(key)
}

func (c *Cache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store.Put(key, value, 0)
}
