// Package array provides a generic slice-based collection implementation that satisfies the abstract.Collection interface.
package array

import (
	"github.com/0x626f/go-kit/pkg/abstract"
	"github.com/0x626f/go-kit/pkg/utils"
)

// Array is a generic slice wrapper that implements the abstract.Collection interface.
type Array[T any, I int] struct {
	items []T
}

// New creates and returns a new empty Array instance.
func New[T any]() *Array[T, int] {
	return &Array[T, int]{}
}

// Wrap creates a new Array containing the provided items.
func Wrap[T any](items ...T) *Array[T, int] {
	instance := New[T]()

	for _, item := range items {
		instance.items = append(instance.items, item)
	}

	return instance
}

// Size returns the number of elements in the array.
func (array *Array[T, I]) Size() int {
	return len(array.items)
}

// IsEmpty returns true if the array contains no elements.
func (array *Array[T, I]) IsEmpty() bool {
	return len(array.items) == 0
}

// At returns the element at the specified index without bounds checking.
func (array *Array[T, I]) At(index int) T {
	return array.items[index]
}

// Get returns the element at the specified index with support for negative indices and wrapping.
// Negative indices count backward from the end, and out-of-bounds indices wrap around.
func (array *Array[T, I]) Get(index int) T {
	n := array.Size()

	if index < 0 {
		index = n + (index % n)
		if index == n {
			index--
		}
	}

	if index >= n {
		index = index % n
	}

	return array.At(index)
}

// Push adds a single element to the end of the array.
func (array *Array[T, I]) Push(item T) {
	array.items = append(array.items, item)
}

// PushAll adds multiple elements to the end of the array.
func (array *Array[T, I]) PushAll(items ...T) {
	array.items = append(array.items, items...)
}

// Join adds all elements from another collection to this array.
func (array *Array[T, I]) Join(collection abstract.Collection[T, int]) {
	for index := range collection.Size() {
		array.Push(collection.At(index))
	}
}

// Merge combines this array with another collection and returns a new array containing elements from both.
func (array *Array[T, I]) Merge(collection abstract.Collection[T, int]) abstract.Collection[T, int] {
	result := New[T]()

	array.ForEach(func(index int, item T) bool {
		result.Push(item)
		return true
	})

	collection.ForEach(func(index int, item T) bool {
		result.Push(item)
		return true
	})

	return result
}

// Delete removes the element at the specified index, without preserving order.
func (array *Array[T, I]) Delete(index int) {
	array.DeleteKeepOrdering(index, false)
}

// DeleteBy removes all elements that satisfy the predicate, without preserving order.
func (array *Array[T, I]) DeleteBy(predicate abstract.Predicate[T]) {
	array.DeleteByKeepOrdering(predicate, false)
}

// DeleteAll removes all elements from the array.
func (array *Array[T, I]) DeleteAll() {
	array.items = nil
}

// DeleteKeepOrdering removes the element at the specified index, with optional order preservation.
// If ordered is true, the original order is preserved but with O(n) complexity.
// If ordered is false, the element is swapped with the last element for O(1) complexity.
func (array *Array[T, I]) DeleteKeepOrdering(index int, ordered bool) {
	lastItemIndex := array.Size() - 1
	if ordered {
		copy(array.items[index:], array.items[index+1:])
		array.items[lastItemIndex] = utils.Zero[T]()
		array.items = array.items[:lastItemIndex]
	} else {
		array.Swap(index, lastItemIndex)
		array.items[lastItemIndex] = utils.Zero[T]()
		array.items = array.items[:lastItemIndex]
	}

}

// DeleteByKeepOrdering removes all elements that satisfy the predicate, with optional order preservation.
func (array *Array[T, I]) DeleteByKeepOrdering(predicate abstract.Predicate[T], ordered bool) {
	for index := 0; index < array.Size(); {
		item := array.At(index)
		if predicate(item) {
			array.DeleteKeepOrdering(index, ordered)
		} else {
			index++
		}
	}
}

// Some returns true if at least one element satisfies the predicate.
func (array *Array[T, I]) Some(predicate abstract.Predicate[T]) bool {

	for _, item := range array.items {
		if predicate(item) {
			return true
		}
	}
	return false
}

// Find returns the first element that satisfies the predicate and a boolean indicating if found.
func (array *Array[T, I]) Find(predicate abstract.Predicate[T]) (T, bool) {
	for _, item := range array.items {
		if predicate(item) {
			return item, true
		}
	}
	return utils.Zero[T](), false
}

// Filter returns a new collection containing only the elements that satisfy the predicate.
func (array *Array[T, I]) Filter(predicate abstract.Predicate[T]) abstract.Collection[T, int] {
	result := New[T]()

	for _, item := range array.items {
		if predicate(item) {
			result.Push(item)
		}
	}

	return result
}

// ForEach executes the provided function once for each element in the collection.
// If the receiver returns false, the iteration stops.
func (array *Array[T, I]) ForEach(receiver abstract.IndexedReceiver[int, T]) {
	for index, item := range array.items {
		if !receiver(index, item) {
			break
		}
	}
}

// First returns the first element of the array.
func (array *Array[T, I]) First() T {
	return array.Get(0)
}

// Last returns the last element of the array.
func (array *Array[T, I]) Last() T {
	return array.Get(-1)
}

// Swap exchanges the elements at the specified indices.
func (array *Array[T, I]) Swap(i, j int) {
	array.items[i], array.items[j] = array.items[j], array.items[i]
}

// Slice creates a new array containing elements from the specified range [from:to).
func (array *Array[T, I]) Slice(from, to int) *Array[T, int] {
	instance := New[T]()
	instance.PushAll(array.items[from:to]...)
	return instance
}

// IndexOf returns the index of the first element that satisfies the predicate,
// or -1 if no element satisfies it.
func (array *Array[T, I]) IndexOf(predicate abstract.Predicate[T]) int {
	for index, item := range array.items {
		if predicate(item) {
			return index
		}
	}
	return -1
}
