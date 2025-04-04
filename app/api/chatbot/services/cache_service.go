package services

import (
	"encoding/json"
	"sync"
	"time"
)

// CacheItem represents an item in the cache with expiration time
type CacheItem struct {
	Value       interface{}
	Expiration  int64         // Unix timestamp when item expires
	Priority    CachePriority // Priority level for cache management
	AccessCount int           // Counter for how often this item is accessed
	LastAccess  int64         // Timestamp of last access
}

// CachePriority defines priority levels for cache items
type CachePriority int

const (
	PriorityLow CachePriority = iota
	PriorityMedium
	PriorityHigh
	PriorityCritical
)

// CacheService provides caching functionality
type CacheService struct {
	items map[string]CacheItem
	mu    sync.RWMutex
	// Default expiration times in seconds for different types of data
	defaultExpirations map[string]int64
	maxItems           int // Maximum items to store before cleanup
}

// NewCacheService creates a new cache service
func NewCacheService() *CacheService {
	cache := &CacheService{
		items: make(map[string]CacheItem),
		defaultExpirations: map[string]int64{
			"drug_info":       86400,   // Drug info: 24 hours
			"trending":        43200,   // Trending data: 12 hours
			"conversation":    3600,    // Conversation context: 1 hour
			"common_searches": 604800,  // Common searches: 1 week
			"side_effects":    172800,  // Side effects: 48 hours
			"recalls":         3600,    // Recalls: 1 hour (refresh more frequently)
			"images":          2592000, // Images: 30 days
		},
		maxItems: 1000, // Set max cache items before cleanup
	}

	// Start janitor routine to clean expired items
	go cache.janitor()

	return cache
}

// Set adds an item to the cache with the specified TTL (time to live) in seconds
func (c *CacheService) Set(key string, value interface{}, ttl int64) {
	c.SetWithPriority(key, value, ttl, PriorityMedium)
}

// SetWithPriority adds an item to the cache with specified TTL and priority
func (c *CacheService) SetWithPriority(key string, value interface{}, ttl int64, priority CachePriority) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var expiration int64
	if ttl > 0 {
		expiration = time.Now().Unix() + ttl
	} else {
		// For zero TTL, use default expiration based on key prefix
		for prefix, defaultTTL := range c.defaultExpirations {
			if len(key) > len(prefix) && key[:len(prefix)] == prefix {
				expiration = time.Now().Unix() + defaultTTL
				break
			}
		}

		// If no matching prefix, default to 1 hour
		if expiration == 0 {
			expiration = time.Now().Unix() + 3600
		}
	}

	// Check if we need to clean up before adding
	if len(c.items) >= c.maxItems {
		c.cleanupLRU()
	}

	c.items[key] = CacheItem{
		Value:       value,
		Expiration:  expiration,
		Priority:    priority,
		AccessCount: 0,
		LastAccess:  time.Now().Unix(),
	}
}

// Get retrieves an item from the cache
func (c *CacheService) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	item, found := c.items[key]
	c.mu.RUnlock()

	if !found {
		return nil, false
	}

	// Check if the item is expired
	if item.Expiration > 0 && item.Expiration < time.Now().Unix() {
		// Item is expired, remove it
		c.mu.Lock()
		delete(c.items, key)
		c.mu.Unlock()
		return nil, false
	}

	// Update access stats
	c.mu.Lock()
	item.AccessCount++
	item.LastAccess = time.Now().Unix()
	c.items[key] = item // Update with new access count
	c.mu.Unlock()

	return item.Value, true
}

// GetJSON retrieves an item from the cache and unmarshals it into the provided value
func (c *CacheService) GetJSON(key string, value interface{}) bool {
	data, found := c.Get(key)
	if !found {
		return false
	}

	// For string data, unmarshal it
	if str, ok := data.(string); ok {
		if err := json.Unmarshal([]byte(str), value); err != nil {
			return false
		}
		return true
	}

	// For already structured data, try direct JSON marshaling then unmarshaling
	if jsonData, err := json.Marshal(data); err == nil {
		if err := json.Unmarshal(jsonData, value); err == nil {
			return true
		}
	}

	// Set value directly if types match
	if jsonValue, ok := data.(interface{}); ok {
		if valuePtr, ok := value.(*interface{}); ok {
			*valuePtr = jsonValue
			return true
		}
	}

	return false
}

// Delete removes an item from the cache
func (c *CacheService) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Flush removes all items from the cache
func (c *CacheService) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]CacheItem)
}

// SetDefaultExpiration sets the default expiration time for a specific type of data
func (c *CacheService) SetDefaultExpiration(dataType string, seconds int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.defaultExpirations[dataType] = seconds
}

// cleanupLRU removes the least recently used items when cache is full
func (c *CacheService) cleanupLRU() {
	if len(c.items) < c.maxItems {
		return
	}

	// First try to remove expired items
	now := time.Now().Unix()
	for k, v := range c.items {
		if v.Expiration < now {
			delete(c.items, k)
		}
	}

	// If still too many items, remove least recently used low priority items
	if len(c.items) >= c.maxItems {
		type itemKey struct {
			key        string
			priority   CachePriority
			lastAccess int64
		}

		var candidates []itemKey

		// Collect low priority items
		for k, v := range c.items {
			if v.Priority <= PriorityMedium {
				candidates = append(candidates, itemKey{k, v.Priority, v.LastAccess})
			}
		}

		// Sort by priority (ascending) and last access time (ascending)
		if len(candidates) > 0 {
			// Simple bubble sort for demonstration
			for i := 0; i < len(candidates); i++ {
				for j := i + 1; j < len(candidates); j++ {
					if candidates[i].priority > candidates[j].priority ||
						(candidates[i].priority == candidates[j].priority &&
							candidates[i].lastAccess > candidates[j].lastAccess) {
						candidates[i], candidates[j] = candidates[j], candidates[i]
					}
				}
			}

			// Remove up to 10% of items
			removeCount := len(c.items) / 10
			if removeCount < 1 {
				removeCount = 1
			}

			for i := 0; i < removeCount && i < len(candidates); i++ {
				delete(c.items, candidates[i].key)
			}
		}
	}
}

// janitor runs periodic cleanup of expired items
func (c *CacheService) janitor() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			now := time.Now().Unix()
			for k, v := range c.items {
				if v.Expiration > 0 && v.Expiration < now {
					delete(c.items, k)
				}
			}
			c.mu.Unlock()
		}
	}
}

// CacheKeys returns all keys in the cache (for debugging)
func (c *CacheService) CacheKeys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]string, 0, len(c.items))
	for k := range c.items {
		keys = append(keys, k)
	}

	return keys
}

// ItemCount returns the number of items in the cache
func (c *CacheService) ItemCount() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

// PreloadCommonDrugs sets cache entries for commonly searched drugs with longer expiration
func (c *CacheService) PreloadCommonDrugs(drugs map[string]interface{}) {
	for drug, data := range drugs {
		c.SetWithPriority("drug_info:"+drug, data, c.defaultExpirations["common_searches"], PriorityHigh)
	}
}
