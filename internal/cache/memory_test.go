// internal/cache/memory_test.go
package cache

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestMemoryCache_BasicOperations(t *testing.T) {
	cache := NewMemoryCache()

	// Test Set and Get
	t.Run("Set and Get", func(t *testing.T) {
		err := cache.Set("key1", "value1", 0)
		if err != nil {
			t.Errorf("Set returned error: %v", err)
		}

		item, err := cache.Get("key1")
		if err != nil {
			t.Errorf("Get returned error: %v", err)
		}
		if item.Value != "value1" {
			t.Errorf("expected value1, got %v", item.Value)
		}
	})

	// Test Get non-existent key
	t.Run("Get non-existent", func(t *testing.T) {
		_, err := cache.Get("nonexistent")
		if err != ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})

	// Test Delete
	t.Run("Delete", func(t *testing.T) {
		cache.Set("key2", "value2", 0)
		err := cache.Delete("key2")
		if err != nil {
			t.Errorf("Delete returned error: %v", err)
		}

		_, err = cache.Get("key2")
		if err != ErrNotFound {
			t.Errorf("expected ErrNotFound after Delete, got %v", err)
		}
	})

	// Test Clear
	t.Run("Clear", func(t *testing.T) {
		cache.Set("key3", "value3", 0)
		cache.Set("key4", "value4", 0)

		err := cache.Clear()
		if err != nil {
			t.Errorf("Clear returned error: %v", err)
		}

		_, err = cache.Get("key3")
		if err != ErrNotFound {
			t.Errorf("expected ErrNotFound after Clear, got %v", err)
		}
	})
}

func TestMemoryCache_Expiration(t *testing.T) {
	cache := NewMemoryCache()

	t.Run("Item expires", func(t *testing.T) {
		// Set item with 100ms expiration
		cache.Set("exp_key", "exp_value", 100*time.Millisecond)

		// Should be able to get it immediately
		item, err := cache.Get("exp_key")
		if err != nil {
			t.Errorf("Get returned error: %v", err)
		}
		if item.Value != "exp_value" {
			t.Errorf("expected exp_value, got %v", item.Value)
		}

		// Wait for expiration
		time.Sleep(150 * time.Millisecond)

		// Should return expired error
		_, err = cache.Get("exp_key")
		if err != ErrExpired {
			t.Errorf("expected ErrExpired, got %v", err)
		}
	})

	t.Run("Zero expiration never expires", func(t *testing.T) {
		cache.Set("never_exp", "value", 0)

		// Wait some time
		time.Sleep(150 * time.Millisecond)

		// Should still be able to get it
		item, err := cache.Get("never_exp")
		if err != nil {
			t.Errorf("Get returned error: %v", err)
		}
		if item.Value != "value" {
			t.Errorf("expected value, got %v", item.Value)
		}
	})
}

func TestMemoryCache_Concurrent(t *testing.T) {
	cache := NewMemoryCache()
	var wg sync.WaitGroup

	// Number of concurrent goroutines
	workers := 10
	// Operations per goroutine
	ops := 100

	t.Run("Concurrent Set and Get", func(t *testing.T) {
		for i := 0; i < workers; i++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()
				for j := 0; j < ops; j++ {
					key := fmt.Sprintf("key_%d_%d", workerID, j)
					value := fmt.Sprintf("value_%d_%d", workerID, j)

					// Set value
					err := cache.Set(key, value, time.Minute)
					if err != nil {
						t.Errorf("Set returned error: %v", err)
					}

					// Get value back
					item, err := cache.Get(key)
					if err != nil {
						t.Errorf("Get returned error: %v", err)
					}
					if item.Value != value {
						t.Errorf("expected %s, got %v", value, item.Value)
					}
				}
			}(i)
		}
		wg.Wait()
	})
}

func TestMemoryCache_Types(t *testing.T) {
	cache := NewMemoryCache()

	t.Run("Different value types", func(t *testing.T) {
		testCases := []struct {
			key   string
			value interface{}
		}{
			{"string", "hello"},
			{"int", 42},
			{"float", 3.14},
			{"bool", true},
			{"struct", struct{ Name string }{"test"}},
			{"slice", []int{1, 2, 3}},
			{"map", map[string]int{"one": 1}},
		}

		for _, tc := range testCases {
			t.Run(tc.key, func(t *testing.T) {
				err := cache.Set(tc.key, tc.value, 0)
				if err != nil {
					t.Errorf("Set returned error: %v", err)
				}

				item, err := cache.Get(tc.key)
				if err != nil {
					t.Errorf("Get returned error: %v", err)
				}

				// Use reflect.DeepEqual for comparing complex types
				if !reflect.DeepEqual(item.Value, tc.value) {
					t.Errorf("expected %v, got %v", tc.value, item.Value)
				}
			})
		}
	})
}

func TestMemoryCache_CleanupExpired(t *testing.T) {
	cache := NewMemoryCache()

	t.Run("Cleanup removes expired items", func(t *testing.T) {
		// Add items with short expiration
		for i := 0; i < 10; i++ {
			key := fmt.Sprintf("key_%d", i)
			cache.Set(key, i, 100*time.Millisecond)
		}

		// Add some non-expiring items
		for i := 0; i < 5; i++ {
			key := fmt.Sprintf("permanent_%d", i)
			cache.Set(key, i, 0)
		}

		// Wait for items to expire and cleanup to run
		time.Sleep(200 * time.Millisecond)

		// Check that expired items are gone
		for i := 0; i < 10; i++ {
			key := fmt.Sprintf("key_%d", i)
			_, err := cache.Get(key)
			if err != ErrExpired && err != ErrNotFound {
				t.Errorf("expected ErrExpired or ErrNotFound for %s, got %v", key, err)
			}
		}

		// Check that non-expiring items remain
		for i := 0; i < 5; i++ {
			key := fmt.Sprintf("permanent_%d", i)
			item, err := cache.Get(key)
			if err != nil {
				t.Errorf("Get returned error for permanent item: %v", err)
			}
			if item.Value != i {
				t.Errorf("expected %d, got %v", i, item.Value)
			}
		}
	})
}
