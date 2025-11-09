package raw

/*
#include <stdlib.h>
#include <string.h>
*/
import "C"
import "unsafe"

func Malloc(size Size) unsafe.Pointer {
	if size == 0 {
		return nil
	}

	ptr := C.malloc(C.size_t(size))

	return ptr
}

func Free(ptr unsafe.Pointer) {
	C.free(ptr)
}

func Void(ptr unsafe.Pointer, size Size) {
	C.memset(ptr, 0, C.size_t(size))
}

func Allocate[T any]() *T {
	return (*T)(Malloc(SizeOf[T]()))
}

func AllocateBlank[T any]() *T {
	size := SizeOf[T]()
	ptr := Malloc(size)

	Void(ptr, size)

	return (*T)(ptr)
}

func DeAllocate[T any](obj *T) {
	Free(unsafe.Pointer(obj))
}

func SizeOf[T any]() Size {
	var sample T
	return unsafe.Sizeof(sample)
}
