package cache

import (
	"sync"
	"time"
)

// Data holds the value and deadline of a key/value pair
type Data struct {
	Value    string
	deadline time.Time
	expired  bool
}

// Cache holds values with type string and
// allows to retrieve them using keys of time strings.
// Key/value pairs can expire if given a deadline using PutTill method.
type Cache struct {
	mu   sync.RWMutex
	data map[string]Data
}

// NewCache returns a new cache instance
func NewCache() Cache {
	return Cache{
		data: map[string]Data{},
	}
}

// Get returns the value associated with the key and the boolean ok (true if exists, false if not),
// if the deadline of the key/value pair has not been exceeded yet.
func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if item, ok := c.data[key]; ok {
		item.expired = item.deadline.Before(time.Now())
		if !item.expired || item.deadline.IsZero() {
			return item.Value, true
		}
	}
	return "", false
}

// Put places a value with an associated key into cache.
// Value put with this method never expired (have infinite deadline).
// Putting into the existing key overwrites the value
func (c *Cache) Put(key, value string) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	c.data[key] = Data{value, time.Time{}, false}
}

// Keys returns a slice of keys in the cache (not expired)
func (c *Cache) Keys() []string {
	c.mu.RLock()
	defer c.mu.RLock()
	keys := make([]string, 0, len(c.data))
	for k, v := range c.data {
		v.expired = v.deadline.Before(time.Now())
		if !v.expired || v.deadline.IsZero() {
			keys = append(keys, k)
		}
	}
	return keys
}

// PutTill places a value with an associated key into cache.
// Value put with this method expires after the given deadline.
// Putting into the existing key overwrites the value
func (c *Cache) PutTill(key, value string, deadline time.Time) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	c.data[key] = Data{value, deadline, false}
}
