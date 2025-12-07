// Package set provides implementations for set data structures that maintain unique elements.
package set

import "github.com/0x626f/go-kit/abstract"

// PrimarySet is a specialized set implementation that works with comparable types directly,
// without requiring them to implement the Keyable interface. It internally wraps values
// in KeyableWrapper to adapt them for use with the standard Set implementation.
//
// Type parameters:
//   - T: The element type, must be comparable (support == and != operators)
type PrimarySet[T comparable] struct {
	*Set[T, *abstract.KeyableWrapper[T]]
}

// NewPrimarySet creates and returns an empty PrimarySet instance.
// This is a convenience wrapper around Set for working with primitive comparable types.
func NewPrimarySet[T comparable]() *PrimarySet[T] {
	return &PrimarySet[T]{
		Set: New[T, *abstract.KeyableWrapper[T]](),
	}
}

// WrapToPrimarySet creates a new PrimarySet containing the provided items.
// This is a convenience function for creating and populating a primary set in one call.
func WrapToPrimarySet[T comparable](items ...T) *PrimarySet[T] {
	instance := NewPrimarySet[T]()
	for _, item := range items {
		instance.Push(instance.Item(item))
	}
	return instance
}

// Item wraps a comparable value in a KeyableWrapper so it can be stored in the set.
// This is a helper method used internally by the PrimarySet implementation.
func (primarySet *PrimarySet[T]) Item(item T) *abstract.KeyableWrapper[T] {
	return &abstract.KeyableWrapper[T]{
		Wrapped: item,
	}
}
