// Package utils provides generic utility functions for common programming tasks.
// This package includes helpers for value manipulation, type conversion, and
// container operations to simplify Go code using generics.
package utils

import "reflect"

// NewInstanceOf creates a new zero-initialized instance of type T and returns a pointer to it.
// This function uses reflection to dynamically create an instance of any type provided as a
// type parameter, which is particularly useful for generic code that needs to instantiate
// types that aren't known until runtime.
//
// Type parameters:
//   - T: Any type to create an instance of
//
// Returns:
//   - *T: A pointer to a new zero-initialized instance of type T
//
// Example:
//
//	type Person struct {
//	    Name string
//	    Age  int
//	}
//
//	// Create a new Person instance
//	p := NewInstanceOf[Person]()
func NewInstanceOf[T any]() *T {
	var sample T
	return reflect.New(reflect.TypeOf(sample)).Interface().(*T)
}

// IsObject determines if the type parameter T is a struct type.
// This is useful for generic code that needs to verify that a type
// is a complex object with fields rather than a primitive type.
//
// Type parameters:
//   - T: Any type to check
//
// Returns:
//   - bool: true if T is a struct type, false otherwise
//
// Example:
//
//	// Returns true
//	isStruct := IsObject[Person]()
//
//	// Returns false
//	isPrimitive := IsObject[int]()
func IsObject[T any]() bool {
	var sample T
	t := reflect.TypeOf(sample)

	return t.Kind() == reflect.Struct
}

// Forward returns a pointer to a copy of the provided value.
// This is useful when you need a pointer to a value but don't want to declare
// a separate variable.
//
// Type parameters:
//   - T: The type of value to forward
//
// Parameters:
//   - arg: The value to create a pointer to
//
// Returns a pointer to a copy of the input value.
func Forward[T any](arg T) *T {
	return &arg
}

// ForwardAll returns a slice of pointers to copies of the provided values.
// This is useful when you need to convert a slice of values to a slice of pointers.
//
// Type parameters:
//   - T: The type of values to forward
//
// Parameters:
//   - args: The values to create pointers to
//
// Returns a slice of pointers to copies of the input values.
func ForwardAll[T any](args ...T) []*T {
	var result []*T
	for _, item := range args {
		result = append(result, &item)
	}
	return result
}

// Zero returns the zero value for the specified type.
// This is useful in generic code when you need a zero value of a type parameter.
//
// Type parameters:
//   - T: The type for which to return a zero value
//
// Returns the zero value of type T.
// For example:
//   - For numeric types: 0
//   - For bool: false
//   - For string: ""
//   - For pointer, channel, slice, map types: nil
//   - For structs: a struct with all fields set to their zero values
func Zero[T any]() T {
	var zero T
	return zero
}

// MapToKeySlice extracts all keys from a map and returns them as a slice.
// This is useful when you need to iterate over map keys in a specific order
// or perform operations on just the keys.
//
// Type parameters:
//   - K: The key type of the map must be comparable
//   - V: The value type of the map
//
// Parameters:
//   - m: The map from which to extract keys
//
// Returns a slice containing all the keys from the input map.
// Note: The order of keys in the returned slice is not deterministic.
func MapToKeySlice[K comparable, V any](m map[K]V) []K {
	slice := make([]K, len(m))

	for key, _ := range m {
		slice = append(slice, key)
	}
	return slice
}

// MapToValueSlice extracts all values from a map and returns them as a slice.
// This is useful when you need to iterate over map values in a specific order
// or perform operations on just the values.
//
// Type parameters:
//   - K: The key type of the map must be comparable
//   - V: The value type of the map
//
// Parameters:
//   - m: The map from which to extract values
//
// Returns a slice containing all the values from the input map.
// Note: The order of values in the returned slice is not deterministic
// and corresponds to the order of keys in the map's internal representation.
func MapToValueSlice[K comparable, V any](m map[K]V) []V {
	slice := make([]V, len(m))

	for _, value := range m {
		slice = append(slice, value)
	}
	return slice
}
