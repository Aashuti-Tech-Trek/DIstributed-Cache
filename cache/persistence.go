package cache

import (
	"encoding/json"
	"os"
	"time"
)

func (c *Cache) SaveToDisk(filename string) error {
	data := make(map[string]entry)
	for k, ele := range c.store.cache {
		data[k] = ele.Value.(entry)
	}

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
		c.store.Put(k, ent.value, ent.expiry-time.Now().Unix())
	}

	return nil
}

func (c *Cache) SetWithTTL(key, value string, ttl int64) {
	c.store.Put(key, value, ttl)
}