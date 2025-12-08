package cache

import (
	"github.com/0x626f/go-kit/linkedlist"
	"github.com/0x626f/go-kit/types"
	"github.com/0x626f/go-kit/utils"
)

// LFUCache implements a Least Frequently Used cache eviction policy.
// When the cache reaches its capacity, it evicts the items that have been accessed
// the least number of times.
//
// The cache tracks access frequency for each item. When an item is accessed via Get,
// its frequency counter is incremented. Items with lower frequency counts are
// evicted first when the cache reaches capacity.
//
// Type parameters:
//   - K: The type of keys (must be comparable)
//   - D: The type of data stored
//
// Time complexity:
//   - Set: O(1)
//   - Get: O(1)
//   - Delete: O(1)
//   - Flush: O(n log n) due to sorting
type LFUCache[K comparable, D any] struct {
	// capacity is the maximum number of unique frequency buckets
	capacity int

	// frequencies is a linked list of frequency buckets
	// Each bucket contains all items with the same access frequency
	frequencies *linkedlist.LinkedList[*types.Pair[uint, PrimaryCache[K, D]]]

	// data maps frequency counts to their corresponding nodes in the frequencies list
	data PrimaryCache[uint, *linkedlist.LinkedNode[*types.Pair[uint, PrimaryCache[K, D]]]]

	// spot maps keys to their frequency bucket nodes for O(1) lookup
	spot PrimaryCache[K, *linkedlist.LinkedNode[*types.Pair[uint, PrimaryCache[K, D]]]]
}

// NewLFUCache creates and initializes a new LFU cache with the specified capacity.
//
// Type parameters:
//   - K: The type of keys (must be comparable)
//   - D: The type of data stored
//
// Parameters:
//   - capacity: Maximum number of frequency buckets the cache can maintain
//
// Returns:
//   - A pointer to the newly created LFUCache
//
// Example:
//
//	cache := cache.NewLFUCache[string, int](100)
//	cache.Set("counter", 1)
//	cache.Get("counter") // Increases frequency
//	cache.Get("counter") // Increases frequency again
func NewLFUCache[K comparable, D any](capacity int) *LFUCache[K, D] {
	return &LFUCache[K, D]{
		capacity:    capacity,
		frequencies: linkedlist.NewLinkedList[*types.Pair[uint, PrimaryCache[K, D]]](),
		data:        make(PrimaryCache[uint, *linkedlist.LinkedNode[*types.Pair[uint, PrimaryCache[K, D]]]]),
		spot:        make(PrimaryCache[K, *linkedlist.LinkedNode[*types.Pair[uint, PrimaryCache[K, D]]]]),
	}
}

// record is an internal method that creates or retrieves a frequency bucket for the given frequency.
// If a bucket for this frequency doesn't exist, it creates one.
//
// Parameters:
//   - freq: The frequency count
//
// Returns:
//   - A pointer to the node containing the frequency bucket
func (cache *LFUCache[K, D]) record(freq uint) *linkedlist.LinkedNode[*types.Pair[uint, PrimaryCache[K, D]]] {
	if cache.data[freq] == nil {
		entry := &types.Pair[uint, PrimaryCache[K, D]]{First: freq, Second: make(PrimaryCache[K, D])}
		node := cache.frequencies.Insert(entry)
		cache.data[freq] = node
	}
	return cache.data[freq]
}

// Set adds a new item to the cache with an initial frequency of 1.
// If the key already exists, this method does nothing (existing value is preserved).
//
// New items are added to the frequency bucket for count 1.
//
// Parameters:
//   - key: The key to associate with the data
//   - item: The data to cache
//
// Time complexity: O(1)
func (cache *LFUCache[K, D]) Set(key K, item D) {
	if _, exists := cache.spot[key]; !exists {
		node := cache.record(1)
		node.Data.Second[key] = item

		cache.spot[key] = node
	}
}

// Get retrieves an item from the cache and increments its access frequency.
// The item is moved to the next frequency bucket (frequency + 1).
//
// Parameters:
//   - key: The key of the item to retrieve
//
// Returns:
//   - The cached data and true if found
//   - A zero value and false if the key is not in the cache
//
// Time complexity: O(1)
func (cache *LFUCache[K, D]) Get(key K) (D, bool) {
	node, exists := cache.spot[key]

	if !exists {
		return utils.Zero[D](), false
	}

	nextNode := cache.record(node.Data.First + 1)

	nextNode.Data.Second[key] = node.Data.Second[key]

	delete(node.Data.Second, key)
	cache.spot[key] = nextNode

	return nextNode.Data.Second[key], true
}

// Delete removes an item from the cache by its key.
//
// Parameters:
//   - key: The key of the item to remove
//
// Returns:
//   - true (always returns true, even if key was not found)
//
// Time complexity: O(1)
func (cache *LFUCache[K, D]) Delete(key K) bool {
	node, exists := cache.spot[key]

	if !exists {
		return true
	}

	delete(node.Data.Second, key)
	delete(cache.spot, key)

	return true
}

// Flush removes frequency buckets when the cache exceeds its capacity.
// It sorts frequency buckets by frequency count (highest first) and keeps
// only the top 'capacity' buckets, removing items in lower frequency buckets.
//
// The flush operation performs the following steps:
//  1. Checks if the number of frequency buckets exceeds capacity
//  2. Sorts frequency buckets by frequency (descending - highest first)
//  3. Keeps only the first 'capacity' buckets
//  4. Removes all items from buckets beyond capacity
//
// This method is useful for periodic cleanup when the number of frequency
// buckets has grown beyond the desired capacity limit.
//
// Time complexity: O(n log n) where n is the number of frequency buckets
func (cache *LFUCache[K, D]) Flush() {
	if cache.frequencies.Size() > cache.capacity {
		cache.frequencies.Sort(func(arg0, arg1 *types.Pair[uint, PrimaryCache[K, D]]) int {
			return int(arg1.First) - int(arg0.First)
		})
		cache.frequencies.ForEach(func(index int, data *types.Pair[uint, PrimaryCache[K, D]]) bool {
			if (index + 1) > cache.capacity {
				for key := range data.Second {
					delete(cache.spot, key)
				}
				delete(cache.data, data.First)
			}
			return true
		})
		cache.frequencies.Shrink(cache.capacity)
	}
}

// Clear removes all items from the cache, resetting it to an empty state.
// This includes clearing all frequency buckets, the frequency-to-node mapping,
// and the key-to-node mapping.
//
// After calling Clear, the cache is empty and ready to accept new items.
// The capacity remains unchanged.
//
// Example:
//
//	cache := NewLFUCache[string, int](100)
//	cache.Set("key1", 1)
//	cache.Set("key2", 2)
//	cache.Clear() // Cache is now empty
//	cache.Set("key3", 3) // Can continue using the cache
//
// Time complexity: O(n + m) where n is the number of items and m is the number of frequency buckets
func (cache *LFUCache[K, D]) Clear() {
	cache.frequencies.DeleteAll()
	clear(cache.data)
	clear(cache.spot)
}
