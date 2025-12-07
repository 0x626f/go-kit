// Package abstract provides generic interfaces for collection types.
package abstract

// Collection is a generic container interface that provides common operations
// for working with collections of elements.
//
// Type parameters:
//   - T: The type of elements stored in the collection
//   - I: The type used for indexing elements
type Collection[I comparable, T any] interface {
	// Size returns the number of elements in the collection.
	Size() int

	// IsEmpty returns true if the collection contains no elements.
	IsEmpty() bool

	// At returns the element at the specified index.
	// This method typically does not perform bounds checking.
	At(I) T

	// Get returns the element at the specified index.
	// This method typically handles negative indices by counting from the end.
	Get(I) T

	// Push adds a single element to the collection.
	Push(T)

	// PushAll adds multiple elements to the collection.
	PushAll(...T)

	// Join adds all elements from another collection to this collection.
	// This operation modifies the current collection.
	Join(Collection[I, T])

	// Merge combines this collection with another collection and returns a new collection.
	// This operation does not modify the current collection.
	Merge(Collection[I, T]) Collection[I, T]

	// Delete removes the element at the specified index.
	Delete(I)

	// DeleteBy removes all elements that satisfy the given predicate.
	DeleteBy(Predicate[T])

	// DeleteAll removes all elements from the collection.
	DeleteAll()

	// Some returns true if at least one element satisfies the predicate.
	Some(Predicate[T]) bool

	// Find returns the first element that satisfies the predicate and a boolean
	// indicating whether such an element was found.
	Find(Predicate[T]) (T, bool)

	// Filter creates a new collection containing only elements that satisfy the predicate.
	Filter(Predicate[T]) Collection[I, T]

	// ForEach executes the provided function once for each element in the collection.
	ForEach(IndexedReceiver[I, T])
}
