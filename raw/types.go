// Package raw provides low-level memory allocation and manipulation functions
// using C's malloc/free and Go's unsafe package.
//
// WARNING: This package uses unsafe operations and direct memory management.
// Improper use can lead to:
//   - Memory leaks
//   - Segmentation faults
//   - Data corruption
//   - Undefined behavior
//
// Only use this package when you need explicit control over memory allocation
// and understand the risks involved. For most use cases, Go's built-in
// memory management is safer and sufficient.
package raw

// Size represents a memory size in bytes, aliased to uintptr for compatibility
// with unsafe.Sizeof and C memory functions.
type Size = uintptr

// Pointer represents a raw memory address, aliased to uintptr.
// This type is used for low-level pointer arithmetic and memory operations.
type Pointer = uintptr
