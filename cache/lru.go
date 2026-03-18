package cache

import (
	"container/list"
	"time"
)

type entry struct {
	key    string
	value  string
	expiry int64
}

type LRUCache struct {
	capacity int
	ll       *list.List
	cache    map[string]*list.Element
}

func NewLRU(cap int) *LRUCache {
	return &LRUCache{
		capacity: cap,
		ll:       list.New(),
		cache:    make(map[string]*list.Element),
	}
}

func (c *LRUCache) Get(key string) (string, bool) {
	if ele, ok := c.cache[key]; ok {
		ent := ele.Value.(entry)
		if time.Now().Unix() > ent.expiry {
			c.ll.Remove(ele)
			delete(c.cache, key)
			return "", false
		}
		c.ll.MoveToFront(ele)
		return ent.value, true
	}
	return "", false
}

func (c *LRUCache) Put(key, value string, ttl int64) {
	exp := time.Now().Unix() + ttl
	if ttl == 0 {
		exp = 9223372036854775807 // max int64
	}
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		ele.Value = entry{key, value, exp}
		return
	}
	ele := c.ll.PushFront(entry{key, value, exp})
	c.cache[key] = ele
	if c.ll.Len() > c.capacity {
		last := c.ll.Back()
		if last != nil {
			c.ll.Remove(last)
			delete(c.cache, last.Value.(entry).key)
		}
	}
}