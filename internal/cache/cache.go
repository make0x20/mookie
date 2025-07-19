// internal/cache/cache.go
package cache

import (
	"errors"
	"time"
)

/*
	Package cache provides a generic caching interface and implementations.

	How to use:
	1. Choose a cache implementation (e.g., MemoryCache)
	2. Create a new cache instance
	3. Use Get/Set/Delete
	4. Handle cache errors appropriately

   Example basic usage:
       cache := cache.NewMemoryCache()

       // Store item with 5 minute expiration
       cache.Set("user:123", userData, 5*time.Minute)

       // Retrieve item
       item, err := cache.Get("user:123")
       if err == nil {
           userData := item.Value.(*UserData)
           // Use userData...
       }

	Example with error handling:
		item, err := cache.Get("key")
		switch err {
		case nil:
			// Use item.Value
		case ErrNotFound:
			// Handle cache miss
		case ErrExpired:
			// Handle expired item
		default:
			// Handle unexpected error
		}

	Notes:
	- Set zero expiration time to disable expiration
	- Expired items should be automatically cleaned up
*/

// Define cache errors
var (
	ErrNotFound = errors.New("cache: key not found")
	ErrExpired  = errors.New("cache: item has expired")
)

// Item represents a cache entry
type Item struct {
	// Value holds the cached data
	Value interface{}
	// ExpiresAt holds the time when this item expires
	// If zero, the item never expires
	ExpiresAt time.Time
}

// Cache defines the interface that cache implementations must satisfy
type Cache interface {
	// Get retrieves an item from the cache by key
	// Returns ErrNotFound if the key doesn't exist
	// Returns ErrExpired if the item has expired
	Get(key string) (*Item, error)

	// Set adds an item to the cache with the specified key and expiration
	// If duration is 0, the item never expires
	// If key already exists, the item will be overwritten
	Set(key string, value interface{}, duration time.Duration) error

	// Delete removes an item from the cache
	// Returns nil if the key was removed or didn't exist
	Delete(key string) error

	// Clear removes all items from the cache
	Clear() error
}
