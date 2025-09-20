package mbta

import (
	"time"
	"sync"
)


// In-Memory Cache

type cacheItem struct {
	value      []Station
	expiration time.Time
}

type Cache struct {
	data map[string]cacheItem
	mux  sync.RWMutex
	ttl  time.Duration
}

var stationCache *Cache

// InitCache initializes the in-memory cache with TTL in seconds
func InitCache(ttlSeconds int) {
	stationCache = &Cache{
		data: make(map[string]cacheItem),
		ttl:  time.Duration(ttlSeconds) * time.Second,
	}
}



// get from cache
func (c *Cache) Get(key string) ([]Station, bool) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	item, exists := c.data[key]
	if !exists || time.Now().After(item.expiration) {
		return nil, false
	}
	return item.value, true
}

// set cache
func (c *Cache) Set(key string, value []Station) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.data[key] = cacheItem{
		value:      value,
		expiration: time.Now().Add(c.ttl),
	}
}