package cache

import (
	"github.com/0x626f/go-kit/linkedlist"
	"github.com/0x626f/go-kit/types"
	"github.com/0x626f/go-kit/utils"
)

// LRUCache implements a Least Recently Used cache eviction policy.
// When the cache reaches its capacity, it evicts the item that was least recently accessed.
//
// The cache maintains access order using a linked list, where the most recently accessed
// items are at the front and least recently accessed items are at the back.
//
// Type parameters:
//   - K: The type of keys (must be comparable)
//   - D: The type of data stored
//
// Time complexity:
//   - Set: O(1)
//   - Get: O(1)
//   - Delete: O(1)
type LRUCache[K comparable, D any] struct {
	// capacity is the maximum number of items the cache can hold
	// A capacity of 0 means unlimited
	capacity int

	// recent is a linked list maintaining items in access order
	// Most recently accessed items are at the front
	recent *linkedlist.LinkedList[*types.Pair[K, D]]

	// data maps keys to their corresponding nodes in the linked list
	// for O(1) lookup and access
	data PrimaryCache[K, *linkedlist.LinkedNode[*types.Pair[K, D]]]
}

// NewLRUCache creates and initializes a new LRU cache with the specified capacity.
//
// Type parameters:
//   - K: The type of keys (must be comparable)
//   - D: The type of data stored
//
// Parameters:
//   - capacity: Maximum number of items the cache can hold. Use 0 for unlimited capacity.
//
// Returns:
//   - A pointer to the newly created LRUCache
//
// Example:
//
//	cache := cache.NewLRUCache[string, int](100)
//	cache.Set("user:123", 42)
//	value, found := cache.Get("user:123")
func NewLRUCache[K comparable, D any](capacity int) *LRUCache[K, D] {
	return &LRUCache[K, D]{
		capacity: capacity,
		recent:   linkedlist.NewLinkedList[*types.Pair[K, D]](),
		data:     make(map[K]*linkedlist.LinkedNode[*types.Pair[K, D]]),
	}
}

// Set adds or updates an item in the cache.
// If the key already exists, this method does nothing (existing value is preserved).
// If the cache is at capacity, the least recently used item is evicted to make room.
//
// The newly added item is placed at the front of the access list (most recently used position).
//
// Parameters:
//   - key: The key to associate with the data
//   - item: The data to cache
//
// Time complexity: O(1)
func (cache *LRUCache[K, D]) Set(key K, item D) {
	if _, exists := cache.data[key]; exists {
		return
	}

	node := cache.recent.InsertFront(&types.Pair[K, D]{First: key, Second: item})
	cache.data[key] = node

	if cache.capacity != 0 && cache.recent.Size() > cache.capacity {
		retired := cache.recent.PopRight()
		delete(cache.data, retired.First)
	}
}

// Get retrieves an item from the cache by its key.
// Accessing an item moves it to the front of the access list (marks it as most recently used).
//
// Parameters:
//   - key: The key of the item to retrieve
//
// Returns:
//   - The cached data and true if found
//   - A zero value and false if the key is not in the cache
//
// Time complexity: O(1)
func (cache *LRUCache[K, D]) Get(key K) (D, bool) {
	if node, exists := cache.data[key]; exists {
		cache.recent.MoveToFront(node)
		return node.Data.Second, true
	}

	return utils.Zero[D](), false
}

// Delete removes an item from the cache by its key.
//
// Parameters:
//   - key: The key of the item to remove
//
// Returns:
//   - true if the item was found and deleted
//   - false if the key was not in the cache
//
// Time complexity: O(1)
func (cache *LRUCache[K, D]) Delete(key K) bool {
	if node, exists := cache.data[key]; exists {
		cache.recent.Remove(node)
		delete(cache.data, key)
		return true
	}
	return false
}

// Flush removes items from the cache when the number of items exceeds capacity.
// It keeps only the most recently accessed items up to the cache's capacity limit.
// Items are removed from the back of the access list (least recently used).
//
// The flush operation performs the following steps:
//  1. Checks if the current size exceeds capacity
//  2. Iterates through items and removes those beyond capacity
//  3. Shrinks the internal list to match capacity
//
// This method is useful for periodic cleanup when items have been added
// without triggering automatic eviction (e.g., when capacity was increased).
//
// Example:
//
//	cache := NewLRUCache[string, int](100)
//	// Add items...
//	cache.Flush() // Ensures cache doesn't exceed 100 items
//
// Time complexity: O(n) where n is the number of items to remove
func (cache *LRUCache[K, D]) Flush() {

	if cache.recent.Size() > cache.capacity {
		cache.recent.ForEach(func(index int, data *types.Pair[K, D]) bool {
			if (index + 1) > cache.capacity {
				delete(cache.data, data.First)
			}
			return true
		})
		cache.recent.Shrink(cache.capacity)
	}
}

// Clear removes all items from the cache, resetting it to an empty state.
// This includes clearing the access order list and the key-to-node mapping.
//
// After calling Clear, the cache is empty and ready to accept new items.
// The capacity remains unchanged.
//
// Example:
//
//	cache := NewLRUCache[string, int](100)
//	cache.Set("key1", 1)
//	cache.Set("key2", 2)
//	cache.Clear() // Cache is now empty
//	cache.Set("key3", 3) // Can continue using the cache
//
// Time complexity: O(n) where n is the number of items in the cache
func (cache *LRUCache[K, D]) Clear() {
	cache.recent.DeleteAll()
	clear(cache.data)
}
