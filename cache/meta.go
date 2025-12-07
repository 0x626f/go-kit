// Package cache provides implementations of various caching strategies
// including LRU (Least Recently Used) and LFU (Least Frequently Used) caches.
//
// Caching helps improve performance by storing frequently accessed data
// in memory for quick retrieval, while automatically evicting less important
// data when capacity limits are reached.
package cache

// Cache defines the interface for a generic cache implementation.
// All cache implementations must support basic operations: setting values,
// retrieving values, and deleting values.
//
// Type parameters:
//   - D: The type of data stored in the cache
//   - K: The type of keys used to identify cached data (must be comparable)
type Cache[D any, K comparable] interface {
	// Set stores a value in the cache with the specified key.
	// If the key already exists, behavior depends on the cache implementation.
	Set(key K, data D)

	// Get retrieves a value from the cache by its key.
	// Returns the cached value and true if found, or a zero value and false if not found.
	// For some cache types (like LRU), accessing a value may affect its eviction priority.
	Get(key K) (D, bool)

	// Delete removes a value from the cache by its key.
	// Returns true if the key was found and deleted, false otherwise.
	Delete(key K) bool
}

// PrimaryCache is a simple map-based cache with no eviction policy.
// It stores key-value pairs without any automatic cleanup or size limits.
//
// Type parameters:
//   - K: The type of keys (must be comparable)
//   - D: The type of data stored
//
// This is a basic building block used internally by more sophisticated caches.
type PrimaryCache[K comparable, D any] map[K]D
