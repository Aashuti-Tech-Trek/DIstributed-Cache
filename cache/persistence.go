package cache

import (
	"encoding/json"
	"os"
	"time"
)

func (c *Cache) SaveToDisk(filename string) error {
	c.mu.RLock()
	data := make(map[string]entry)
	for k, ele := range c.store.cache {
		data[k] = ele.Value.(entry)
	}
	c.mu.RUnlock()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(data)
}

func (c *Cache) LoadFromDisk(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	data := make(map[string]entry)
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return err
	}

	for k, ent := range data {
		ttl := ent.Expiry - time.Now().Unix()
		if ttl < 0 && ent.Expiry != 9223372036854775807 {
			continue
		}
		if ent.Expiry == 9223372036854775807 {
			ttl = 0 // Will be handled correctly by Put
		}
		c.SetWithTTL(k, ent.Value, ttl)
	}

	return nil
}

func (c *Cache) SetWithTTL(key, value string, ttl int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store.Put(key, value, ttl)
}
