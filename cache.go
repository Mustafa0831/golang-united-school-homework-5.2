package cache

import (
	"sync"
	"time"
)

type Data struct {
	Value    string
	deadline *time.Time
}
type Cache struct {
	mu   sync.Mutex
	data map[string]Data
}

func NewCache() Cache {
	return Cache{
		data: map[string]Data{},
	}
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.data[key]; !ok {
		return "", false
	}
	if c.data[key].deadline != nil && c.data[key].deadline.Before(time.Now()) {
		return "", false
	}
	return c.data[key].Value, true
}

func (c *Cache) Put(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = Data{value, nil}
}

func (c *Cache) Keys() []string {
	c.mu.Lock()
	defer c.mu.Lock()
	var key []string
	now := time.Now()
	for k, v := range c.data {
		if v.deadline != nil && v.deadline.Before(now) {
			continue
		}
		key = append(key, k)
	}
	return key
}

func (c *Cache) PutTill(key, value string, deadline time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = Data{value, &deadline}
}
