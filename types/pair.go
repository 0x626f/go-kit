// Package types provides common utility types and structures that are shared
// across different packages in the go-kit library.
package types

// Pair is a generic container that holds two values of potentially different types.
// It's useful for functions that need to return two related values or for storing
// key-value associations.
//
// Type parameters:
//   - F: The type of the first element
//   - S: The type of the second element
//
// Example:
//
//	// Create a pair of string and int
//	p := Pair[string, int]{First: "age", Second: 25}
//	fmt.Println(p.First, p.Second) // Output: age 25
//
//	// Create a pair of two strings
//	coordinate := Pair[float64, float64]{First: 10.5, Second: 20.3}
type Pair[F, S any] struct {
	// First holds the first value of the pair
	First F
	// Second holds the second value of the pair
	Second S
}
