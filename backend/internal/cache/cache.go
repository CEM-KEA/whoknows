package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

// wrapper struct for go-cache library Cache
type Cache struct {
	cache *cache.Cache
}

// creates a new Cache instance
func NewCache() *Cache {
	return &Cache{
		cache: cache.New(cache.NoExpiration, cache.NoExpiration),
	}
}

// adds a key-value pair to the cache, with an expiration time
func (c *Cache) Set(key string, value interface{}, expiration time.Duration) {
	c.cache.Set(key, value, expiration)
}

// retrieves a value from the cache
func (c *Cache) Get(key string) (interface{}, bool) {
	return c.cache.Get(key)
}
