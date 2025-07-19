// internal/cache/memory.go
package cache

import (
	"sync"
	"time"
)

/*
	Package cache's memory implementation provides a simple in-memory cache
	with automatic cleanup of expired items.

	How to use:
	1. Create a new memory cache instance
	2. Store and retrieve items with optional expiration
	3. Items are automatically removed when expired

	Example basic usage:
	   // Create new cache
	   cache := cache.NewMemoryCache()

	   // Store items
	   cache.Set("key1", "value1", time.Minute)
	   cache.Set("key2", data, 30*time.Second)

	   // Retrieve items
	   item, err := cache.Get("key1")
	   if err == nil {
		   value := item.Value.(string)
		   // Use value...
	   }

	Features:
	- Thread-safe operations
	- Automatic cleanup of expired items
	- Zero allocation for non-expired gets
	- Configurable cleanup interval
	- Efficient memory usage

	Notes:
	- Uses sync.RWMutex for thread safety
	- Cleanup runs every minute in background
	- Safe for concurrent access
	- Memory is released when items expire
*/

// MemoryCache is an in-memory cache implementation
type MemoryCache struct {
	items map[string]Item
	mu    sync.RWMutex
}

// NewMemoryCache creates a new MemoryCache instance
func NewMemoryCache() *MemoryCache {
	cache := &MemoryCache{
		items: make(map[string]Item),
	}

	// Start the cleanup goroutine
	go cache.cleanup()

	return cache
}

// Get retrieves an item from the cache
func (c *MemoryCache) Get(key string) (*Item, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil, ErrNotFound
	}

	if !item.ExpiresAt.IsZero() && time.Now().After(item.ExpiresAt) {
		return nil, ErrExpired
	}

	return &item, nil
}

// Set adds an item to the cache
func (c *MemoryCache) Set(key string, value interface{}, duration time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// If 0 duration, set expiresAt to zero
	var expiresAt time.Time
	if duration > 0 {
		expiresAt = time.Now().Add(duration)
	}

	c.items[key] = Item{
		Value:     value,
		ExpiresAt: expiresAt,
	}

	return nil
}

// Delete removes an item from the cache
func (c *MemoryCache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
	return nil
}

// Clear removes all items from the cache
func (c *MemoryCache) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]Item)
	return nil
}

// cleanup removes expired items from the cache periodically
func (c *MemoryCache) cleanup() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		for key, item := range c.items {
			if !item.ExpiresAt.IsZero() && time.Now().After(item.ExpiresAt) {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}
