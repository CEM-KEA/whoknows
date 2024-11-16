package cache

import (
	"time"

	"github.com/CEM-KEA/whoknows/backend/internal/utils"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

// Cache is a wrapper struct for the go-cache library.
type Cache struct {
	cache *cache.Cache
}

// NewCache creates a new Cache instance.
func NewCache() *Cache {
	return &Cache{
		cache: cache.New(cache.NoExpiration, cache.NoExpiration),
	}
}

// Set adds a key-value pair to the cache, with an expiration time.
func (c *Cache) Set(key string, value interface{}, expiration time.Duration) {
	utils.LogInfo("Adding item to cache", logrus.Fields{
		"key":        key,
		"expiration": expiration,
	})
	c.cache.Set(key, value, expiration)
}

// Get retrieves a value from the cache and updates Prometheus metrics.
func (c *Cache) Get(key string) (interface{}, bool) {
	utils.LogInfo("Fetching item from cache", logrus.Fields{
		"key": key,
	})
	value, found := c.cache.Get(key)

	// Update Prometheus metrics for cache hits and misses
	cacheName := "default" // Replace with a meaningful cache name if needed
	if found {
		utils.IncrementCacheHit(cacheName)
		utils.LogInfo("Cache hit", logrus.Fields{
			"key": key,
		})
	} else {
		utils.IncrementCacheMiss(cacheName)
		utils.LogInfo("Cache miss", logrus.Fields{
			"key": key,
		})
	}

	return value, found
}
