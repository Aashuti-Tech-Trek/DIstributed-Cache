package cache

type Cache struct {
	store *LRUCache
}

func NewCache(capacity int) *Cache {
	return &Cache{
		store: NewLRU(capacity),
	}
}

func (c *Cache) Get(key string) (string, bool) {
	return c.store.Get(key)
}

func (c *Cache) Set(key, value string) {
	c.store.Put(key, value, 0)
}