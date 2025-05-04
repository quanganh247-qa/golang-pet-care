package services

import (
	"encoding/json"
	"sync"
	"time"
)

// CacheItem represents an item in the cache with expiration
type CacheItem struct {
	Value      interface{}
	Expiration int64
}

// CacheService provides in-memory caching functionality
type CacheService struct {
	items map[string]CacheItem
	mu    sync.RWMutex
}

// NewCacheService creates a new cache service
func NewCacheService() *CacheService {
	cache := &CacheService{
		items: make(map[string]CacheItem),
	}

	// Start a goroutine to clean up expired items
	go cache.cleanupLoop()

	return cache
}

// cleanupLoop runs periodically to remove expired items
func (c *CacheService) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.DeleteExpired()
		}
	}
}

// Set adds a key-value pair to the cache with an expiration time
func (c *CacheService) Set(key string, value interface{}, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var expiration int64
	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}

	c.items[key] = CacheItem{
		Value:      value,
		Expiration: expiration,
	}
}

// Get retrieves a value from the cache
func (c *CacheService) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		return nil, false
	}

	// Check if item has expired
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

// DeleteExpired removes all expired items from the cache
func (c *CacheService) DeleteExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now().UnixNano()

	for key, item := range c.items {
		if item.Expiration > 0 && item.Expiration < now {
			delete(c.items, key)
		}
	}
}

// CacheKeys returns all keys in the cache
func (c *CacheService) CacheKeys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]string, 0, len(c.items))
	for k := range c.items {
		keys = append(keys, k)
	}

	return keys
}

// SetJSON marshals a value to JSON and stores it
func (c *CacheService) SetJSON(key string, value interface{}, duration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	c.Set(key, data, duration)
	return nil
}

// GetJSON retrieves a JSON-serialized value and unmarshals it
func (c *CacheService) GetJSON(key string, target interface{}) bool {
	value, found := c.Get(key)
	if !found {
		return false
	}

	data, ok := value.([]byte)
	if !ok {
		return false
	}

	err := json.Unmarshal(data, target)
	if err != nil {
		return false
	}

	return true
}
