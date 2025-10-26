// Package set provides implementations for set data structures that maintain unique elements.
package set

import (
	"github.com/0x626f/go-kit/abstract"
	"github.com/0x626f/go-kit/utils"
)

// Set represents a collection of unique elements that implement the Keyable interface.
// Elements are stored internally in a map with their keys as the map keys.
//
// Type parameters:
//   - T: The type of elements stored in the set, must implement abstract.Keyable[I]
//   - I: The key type, must be comparable (support == and != operators)
type Set[T abstract.Keyable[I], I comparable] struct {
	data map[I]T
}

// New creates and returns an empty Set instance.
func New[T abstract.Keyable[I], I comparable]() *Set[T, I] {
	return &Set[T, I]{data: make(map[I]T)}
}

// Wrap creates a new Set containing the provided items.
// This is a convenience function for creating and populating a set in one call.
func Wrap[T abstract.Keyable[I], I comparable](items ...T) *Set[T, I] {
	instance := New[T, I]()
	instance.PushAll(items...)
	return instance
}

// Size returns the number of elements in the set.
// Implements the Size method of the abstract.Collection interface.
func (set *Set[T, I]) Size() int {
	return len(set.data)
}

// IsEmpty returns true if the set contains no elements.
// Implements the IsEmpty method of the abstract.Collection interface.
func (set *Set[T, I]) IsEmpty() bool {
	return len(set.data) == 0
}

// At returns the element with the specified key.
// If the key doesn't exist, it returns the zero value of type T.
// Implements the At method of the abstract.Collection interface.
func (set *Set[T, I]) At(key I) T {
	return set.data[key]
}

// Get returns the element with the specified key.
// If the key doesn't exist, it returns the zero value of type T.
// Implements the Get method of the abstract.Collection interface.
func (set *Set[T, I]) Get(key I) T {
	return set.data[key]
}

// Push adds an item to the set if it's not already present.
// Uniqueness is determined by the item's Key() value.
// Implements the Push method of the abstract.Collection interface.
func (set *Set[T, I]) Push(item T) {
	if !set.Has(item) {
		key := item.Key()
		set.data[key] = item
	}
}

// PushAll adds multiple items to the set.
// Only items with unique keys will be added.
// Implements the PushAll method of the abstract.Collection interface.
func (set *Set[T, I]) PushAll(items ...T) {
	for _, item := range items {
		set.Push(item)
	}
}

// Join adds all elements from another collection to this set.
// Implements the Join method of the abstract.Collection interface.
func (set *Set[T, I]) Join(collection abstract.Collection[T, I]) {
	collection.ForEach(func(index I, item T) bool {
		set.Push(item)
		return true
	})
}

// Merge combines this set with another collection and returns a new set.
// This operation doesn't modify the original set.
// Implements the Merge method of the abstract.Collection interface.
func (set *Set[T, I]) Merge(collection abstract.Collection[T, I]) abstract.Collection[T, I] {
	result := New[T, I]()

	set.ForEach(func(index I, item T) bool {
		result.Push(item)
		return true
	})

	collection.ForEach(func(index I, item T) bool {
		result.Push(item)
		return true
	})

	return result
}

// Delete removes the element with the specified key.
// If the key doesn't exist, this is a no-op.
// Implements the Delete method of the abstract.Collection interface.
func (set *Set[T, I]) Delete(index I) {
	delete(set.data, index)
}

// DeleteBy removes all elements that satisfy the predicate.
// Implements the DeleteBy method of the abstract.Collection interface.
func (set *Set[T, I]) DeleteBy(predicate abstract.Predicate[T]) {
	for index, item := range set.data {
		if predicate(item) {
			set.Delete(index)
		}
	}
}

// DeleteAll removes all elements from the set.
// Implements the DeleteAll method of the abstract.Collection interface.
func (set *Set[T, I]) DeleteAll() {
	set.data = make(map[I]T)
}

// Some returns true if at least one element satisfies the predicate.
// Returns false if the set is empty or no element satisfies the predicate.
// Implements the Some method of the abstract.Collection interface.
func (set *Set[T, I]) Some(predicate abstract.Predicate[T]) bool {
	for _, item := range set.data {
		if predicate(item) {
			return true
		}
	}
	return false
}

// Find returns the first element that satisfies the predicate and a boolean
// indicating whether such element was found.
// Implements the Find method of the abstract.Collection interface.
func (set *Set[T, I]) Find(predicate abstract.Predicate[T]) (T, bool) {
	for _, item := range set.data {
		if predicate(item) {
			return item, true
		}
	}

	return utils.Zero[T](), false
}

// Has checks if the set contains the specified item.
// Returns true if the item exists in the set, false otherwise.
func (set *Set[T, I]) Has(item T) bool {
	_, has := set.data[item.Key()]
	return has
}

// Filter creates a new set containing only elements that satisfy the predicate.
// Implements the Filter method of the abstract.Collection interface.
func (set *Set[T, I]) Filter(predicate abstract.Predicate[T]) abstract.Collection[T, I] {
	result := New[T, I]()

	for _, item := range set.data {
		if predicate(item) {
			result.Push(item)
		}
	}

	return result
}

// ForEach executes the provided function once for each element in the set.
// The iteration order is not guaranteed as sets are unordered collections.
// If the receiver function returns false, the iteration is stopped.
// Implements the ForEach method of the abstract.Collection interface.
func (set *Set[T, I]) ForEach(receiver abstract.IndexedReceiver[I, T]) {
	for index, item := range set.data {
		if !receiver(index, item) {
			break
		}
	}
}
