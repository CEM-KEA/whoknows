package cache

import (
	"time"

	"github.com/CEM-KEA/whoknows/backend/internal/utils"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

type Cache struct {
	cache *cache.Cache
}


// NewCache creates and returns a new instance of Cache with no expiration settings.
// It initializes the internal cache with no expiration for both items and cleanup intervals.
func NewCache() *Cache {
	return &Cache{
		cache: cache.New(cache.NoExpiration, cache.NoExpiration),
	}
}


// Set adds an item to the cache with the specified key, value, and expiration duration.
// It logs the action with the key and expiration details.
//
// Parameters:
//   - key: The key under which the value will be stored.
//   - value: The value to be stored in the cache.
//   - expiration: The duration for which the item should remain in the cache.
func (c *Cache) Set(key string, value interface{}, expiration time.Duration) {
	utils.LogInfo("Adding item to cache", logrus.Fields{
		"key":        key,
		"expiration": expiration,
	})
	c.cache.Set(key, value, expiration)
}


// Get retrieves an item from the cache using the provided key.
// It logs the operation and updates Prometheus metrics for cache hits and misses.
//
// Parameters:
//   - key: The key of the item to retrieve from the cache.
//
// Returns:
//   - value: The value associated with the key, if found.
//   - found: A boolean indicating whether the item was found in the cache.
func (c *Cache) Get(key string) (interface{}, bool) {
	utils.LogInfo("Fetching item from cache", logrus.Fields{
		"key": key,
	})
	value, found := c.cache.Get(key)

	// Update Prometheus metrics for cache hits and misses
	cacheName := "default"
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
