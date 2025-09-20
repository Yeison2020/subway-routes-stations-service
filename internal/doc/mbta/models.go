package mbta

import (
	"sync"
	"time"
)

// Models

// ---------------------------
// Models
// ---------------------------

type Station struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Route struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Stations []Station `json:"stations"`
}

// MBTA API responses
type routesResponse struct {
	Data []struct {
		ID         string `json:"id"`
		Attributes struct {
			Name string `json:"long_name"`
		} `json:"attributes"`
	} `json:"data"`
}

type stopsResponse struct {
	Data []struct {
		ID         string `json:"id"`
		Attributes struct {
			Name string `json:"name"`
		} `json:"attributes"`
	} `json:"data"`
}

// ---------------------------
// In-Memory Cache
// ---------------------------

type cacheItem struct {
	value      []Station
	expiration time.Time
}

type Cache struct {
	data map[string]cacheItem
	mux  sync.RWMutex
	ttl  time.Duration
}


