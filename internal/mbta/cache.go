package mbta

import (
	"time"
	"sync"
)

// In Memory Cache 

type cacheItem struct {

	value []Station

	experation time.Time
}


type Cache struct {
	data map[string]cacheItem
	mux sync.RWMutex
	ttl time.Duration
}




var StationCache *Cache

// InitCache initializes the in-memory cache with TTL in seconds

func InitCache(ttlsSeconds int){
	StationCache = &Cache{
		data : make(map[string]cacheItem),
		ttl: time.Duration(ttlsSeconds) * time.Second,
	}
}

// Get from Cache
// Mux to read without other readers

func (c *Cache) Get(key string) ([]Station, bool) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	item, exits := c.data[key]

	if !exits || time.Now().After(item.experation){
		return nil, false
	}
	return item.value, true
}

// set Cache

func (c *Cache) Set(key string, value []Station){
	c.mux.Lock()
	defer c.mux.Unlock()
	c.data[key] = cacheItem{
		value: value,
		experation: time.Now().Add(c.ttl),
	}
}


