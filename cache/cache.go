// Package cache provides implementations of various caching strategies
// including LRU (Least Recently Used) and LFU (Least Frequently Used) caches.
//
// Caching helps improve performance by storing frequently accessed data
// in memory for quick retrieval, while automatically evicting less important
// data when capacity limits are reached.
package cache

// Cache defines the interface for a generic cache implementation.
// All cache implementations must support basic operations: setting values,
// retrieving values, deleting values, clearing all data, and flushing to capacity.
//
// The interface is designed to work with various eviction strategies such as
// LRU (Least Recently Used) and LFU (Least Frequently Used).
//
// Type parameters:
//   - D: The type of data stored in the cache
//   - K: The type of keys used to identify cached data (must be comparable)
//
// Thread Safety:
//
// Cache implementations are not thread-safe by design. For concurrent access,
// wrap the cache with appropriate synchronization primitives (e.g., sync.RWMutex).
type Cache[D any, K comparable] interface {
	// Set stores a value in the cache with the specified key.
	// If the key already exists, the existing value is preserved (not updated).
	// When the cache is at capacity, behavior varies by implementation:
	//   - LRU: Evicts the least recently used item
	//   - LFU: Items are organized by frequency; call Flush to manage capacity
	Set(key K, data D)

	// Get retrieves a value from the cache by its key.
	// Returns the cached value and true if found, or a zero value and false if not found.
	//
	// Side effects vary by implementation:
	//   - LRU: Marks the item as most recently used
	//   - LFU: Increments the item's access frequency
	Get(key K) (D, bool)

	// Delete removes a value from the cache by its key.
	// Returns true if the key was found and deleted, false otherwise.
	//
	// Implementation-specific behavior:
	//   - LRU: Always returns true (for historical reasons)
	//   - LFU: Always returns true (for historical reasons)
	Delete(key K) bool

	// Clear removes all items from the cache, resetting it to an empty state.
	// The cache's capacity remains unchanged and the cache can be reused
	// immediately after calling Clear.
	//
	// This operation is useful for:
	//   - Invalidating all cached data at once
	//   - Resetting cache state during testing
	//   - Responding to memory pressure
	Clear()

	// Flush enforces capacity constraints by removing items based on the
	// cache's eviction policy. The exact behavior depends on the implementation:
	//
	//   - LRU: Removes items beyond capacity, keeping only the most recently used
	//   - LFU: Removes low-frequency buckets when bucket count exceeds capacity
	//
	// This operation is useful for:
	//   - Periodic cleanup to enforce capacity limits
	//   - Reclaiming memory when the cache has grown beyond desired size
	//   - Ensuring consistent cache size after bulk operations
	Flush()
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
