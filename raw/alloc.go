package raw

/*
#include <stdlib.h>
#include <string.h>
*/
import "C"
import "unsafe"

// Malloc allocates size bytes of uninitialized memory using C's malloc function.
// The allocated memory must be freed using Free to avoid memory leaks.
//
// Parameters:
//   - size: The number of bytes to allocate
//
// Returns:
//   - An unsafe.Pointer to the allocated memory, or nil if size is 0
//
// WARNING: The caller is responsible for freeing the allocated memory using Free.
// Failure to do so will result in memory leaks.
//
// Example:
//
//	size := raw.Size(1024)
//	ptr := raw.Malloc(size)
//	defer raw.Free(ptr)
//	// Use ptr...
func Malloc(size Size) unsafe.Pointer {
	if size == 0 {
		return nil
	}

	ptr := C.malloc(C.size_t(size))

	return ptr
}

// Free deallocates memory previously allocated by Malloc using C's free function.
//
// Parameters:
//   - ptr: The pointer to the memory to be freed (must have been allocated by Malloc)
//
// WARNING: Calling Free on a pointer not returned by Malloc, or calling Free
// twice on the same pointer, results in undefined behavior and will likely crash.
func Free(ptr unsafe.Pointer) {
	C.free(ptr)
}

// Void zeros out a block of memory using C's memset function.
// This sets all bytes in the specified memory region to zero.
//
// Parameters:
//   - ptr: Pointer to the memory region to zero out
//   - size: Number of bytes to zero
//
// Example:
//
//	ptr := raw.Malloc(100)
//	raw.Void(ptr, 100) // Set all 100 bytes to zero
//	defer raw.Free(ptr)
func Void(ptr unsafe.Pointer, size Size) {
	C.memset(ptr, 0, C.size_t(size))
}

// Allocate allocates memory for a value of type T and returns a typed pointer.
// The memory is NOT initialized and may contain garbage values.
//
// Type parameters:
//   - T: The type to allocate memory for
//
// Returns:
//   - A pointer to the allocated memory of type *T
//
// WARNING: The allocated memory must be freed using DeAllocate.
// The memory is uninitialized and may contain random data.
//
// Example:
//
//	type MyStruct struct { x, y int }
//	ptr := raw.Allocate[MyStruct]()
//	defer raw.DeAllocate(ptr)
//	ptr.x = 10
//	ptr.y = 20
func Allocate[T any]() *T {
	return (*T)(Malloc(SizeOf[T]()))
}

// AllocateBlank allocates memory for a value of type T and initializes it to zero.
// Unlike Allocate, this function guarantees the memory is zero-initialized.
//
// Type parameters:
//   - T: The type to allocate memory for
//
// Returns:
//   - A pointer to the zero-initialized memory of type *T
//
// WARNING: The allocated memory must be freed using DeAllocate.
//
// Example:
//
//	type MyStruct struct { x, y int }
//	ptr := raw.AllocateBlank[MyStruct]()
//	defer raw.DeAllocate(ptr)
//	// ptr.x and ptr.y are guaranteed to be 0
func AllocateBlank[T any]() *T {
	size := SizeOf[T]()
	ptr := Malloc(size)

	Void(ptr, size)

	return (*T)(ptr)
}

// DeAllocate frees memory previously allocated by Allocate or AllocateBlank.
//
// Type parameters:
//   - T: The type of the allocated object
//
// Parameters:
//   - obj: Pointer to the object to deallocate (must have been allocated by Allocate or AllocateBlank)
//
// WARNING: Do not use obj after calling DeAllocate. Accessing freed memory
// results in undefined behavior.
func DeAllocate[T any](obj *T) {
	Free(unsafe.Pointer(obj))
}

// SizeOf returns the size in bytes of a value of type T.
// This is equivalent to unsafe.Sizeof but works with generic types.
//
// Type parameters:
//   - T: The type to get the size of
//
// Returns:
//   - The size in bytes of type T
//
// Example:
//
//	size := raw.SizeOf[int]()        // Returns 8 on 64-bit systems
//	size := raw.SizeOf[byte]()       // Returns 1
//	size := raw.SizeOf[struct{}]()   // Returns 0
func SizeOf[T any]() Size {
	var sample T
	return unsafe.Sizeof(sample)
}
