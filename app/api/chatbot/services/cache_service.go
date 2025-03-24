package services

import (
	"encoding/json"
	"sync"
	"time"
)

// CacheItem represents a single cached item with expiration time
type CacheItem struct {
	Value      interface{}
	Expiration int64
}

// CacheService provides caching functionality
type CacheService struct {
	items map[string]CacheItem
	mu    sync.RWMutex
}

// NewCacheService creates a new cache service instance
func NewCacheService() *CacheService {
	cache := &CacheService{
		items: make(map[string]CacheItem),
	}
	
	// Start a background goroutine to clean expired items
	go cache.startJanitor()
	
	return cache
}

// Set adds an item to the cache with a TTL (time-to-live) in seconds
func (c *CacheService) Set(key string, value interface{}, ttl int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	var expiration int64
	if ttl > 0 {
		expiration = time.Now().Add(time.Duration(ttl) * time.Second).UnixNano()
	}
	
	c.items[key] = CacheItem{
		Value:      value,
		Expiration: expiration,
	}
}

// Get retrieves an item from the cache
func (c *CacheService) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	item, found := c.items[key]
	if !found {
		return nil, false
	}
	
	// Check if item is expired
	if item.Expiration > 0 && item.Expiration < time.Now().UnixNano() {
		return nil, false
	}
	
	return item.Value, true
}

// Delete removes an item from the cache
func (c *CacheService) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	delete(c.items, key)
}

// Clear empties the entire cache
func (c *CacheService) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.items = make(map[string]CacheItem)
}

// GetTTL returns the remaining time to live for a cached item in seconds
func (c *CacheService) GetTTL(key string) int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	item, found := c.items[key]
	if !found || item.Expiration == 0 {
		return -1
	}
	
	ttl := item.Expiration - time.Now().UnixNano()
	if ttl <= 0 {
		return -1
	}
	
	return ttl / int64(time.Second)
}

// startJanitor starts a background routine to clean expired items
func (c *CacheService) startJanitor() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			c.cleanExpired()
		}
	}
}

// cleanExpired removes expired items from the cache
func (c *CacheService) cleanExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	now := time.Now().UnixNano()
	for k, v := range c.items {
		if v.Expiration > 0 && v.Expiration < now {
			delete(c.items, k)
		}
	}
}

// GetJSON retrieves a JSON-encodable value from the cache and decodes it into the target variable
func (c *CacheService) GetJSON(key string, target interface{}) bool {
	val, found := c.Get(key)
	if !found {
		return false
	}
	
	// Convert to JSON string first
	jsonBytes, err := json.Marshal(val)
	if err != nil {
		return false
	}
	
	// Unmarshal into the target
	err = json.Unmarshal(jsonBytes, target)
	return err == nil
}