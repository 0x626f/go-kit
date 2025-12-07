// Package number provides arbitrary-precision numeric types and utilities.
package number

import "github.com/0x626f/go-kit/abstract"

// IntegerField is a constraint interface that matches all built-in integer types.
// This allows generic functions to operate on any integer type.
//
// Includes both signed (~int, ~int8, etc.) and unsigned (~uint, ~uint8, etc.) integers.
// The ~ operator allows for type aliases as well as the base types.
type IntegerField interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~int | ~int8 | ~int16 | ~int32 | ~int64
}

// RealField is a constraint interface that matches all built-in floating-point types.
// This allows generic functions to operate on any floating-point type.
//
// Includes both single precision (~float32) and double precision (~float64) floating-point numbers.
// The ~ operator allows for type aliases as well as the base types.
type RealField interface {
	~float32 | ~float64
}

// NumericField is a constraint interface that matches all built-in numeric types.
// This allows generic functions to operate on any numeric primitive type.
//
// Includes both integer types (through IntegerField) and floating-point types (through RealField).
type NumericField interface {
	IntegerField | RealField
}

// BigNumericField is an interface that defines operations common to all arbitrary-precision
// numeric types in this package (BigInt and BigFloat).
//
// Type parameters:
//   - T: The concrete numeric type (BigInt or BigFloat)
//
// It extends the abstract.Comparable interface to ensure all big number types
// can be compared with each other.
type BigNumericField[T BigInt | BigFloat] interface {
	// Embed the Comparable interface to require comparison capabilities
	abstract.Comparable[T]

	// IsMutable returns whether the number can be modified by operations.
	IsMutable() bool

	// Basic arithmetic operations
	// Each returns a new number or modifies the receiver if mutable
	Add(*T) *T
	Subtract(*T) *T
	Multiply(*T) *T
	Divide(*T) *T

	// Mathematical functions
	Sqrt() *T
	Abs() *T
	Negate() *T

	// Sign returns the sign of the number (-1, 0, or 1)
	Sign() int
}

// ForwardToBigNumeric converts a pointer to a concrete type (BigInt or BigFloat)
// to its BigNumericField interface representation.
//
// This is a helper function for working with the BigNumericField interface
// when you have a concrete type instance.
//
// Type parameters:
//   - T: The target interface type (must implement BigNumericField[B])
//   - B: The concrete type (BigInt or BigFloat)
//
// Returns the same instance but typed as the BigNumericField interface.
func ForwardToBigNumeric[T BigNumericField[B], B BigInt | BigFloat](arg *B) T {
	return any(arg).(T)
}
