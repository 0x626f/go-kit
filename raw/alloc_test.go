package raw

import (
	"testing"
	"unsafe"
)

func TestMallocAndFree(t *testing.T) {
	var bytes Size = 20

	ptr := Malloc(bytes)
	defer Free(ptr)

	t.Logf("Allocated: %p\n", ptr)
	if ptr == nil {
		t.Fatal("the memory was not allocated")
	}
}

func TestSizeOf(t *testing.T) {

	type stub struct {
		one   int
		two   uint
		three float64
	}

	var one int
	var two uint
	var three float64

	expectedSize := unsafe.Sizeof(one) + unsafe.Sizeof(two) + unsafe.Sizeof(three)
	structSize := SizeOf[stub]()

	if structSize != expectedSize {
		t.Fatal("not matched size")
	}

}

func TestGenericSizeOf(t *testing.T) {

	type stub[T any] struct {
		data T
	}

	type underlying struct {
		name    string
		payload []byte
	}

	var sample stub[underlying]

	expectedSize := unsafe.Sizeof(sample)
	structSize := SizeOf[stub[underlying]]()

	if structSize != expectedSize {
		t.Fatal("not matched size")
	}

}

func TestAllocateAndFree(t *testing.T) {

	type stub struct {
		one   int
		two   uint
		three float64
	}

	ptr := Allocate[stub]()
	defer DeAllocate(ptr)

	t.Logf("Allocated: %p\n", ptr)

	ptr.one = 1
	ptr.two = 1
	ptr.three = 1

}

func BenchmarkMalloc(b *testing.B) {
	b.ReportAllocs()
	sizes := []Size{8, 16, 32, 64, 128, 256, 512, 1024}

	for _, size := range sizes {
		b.Run(string(rune('0'+size/8)), func(b *testing.B) {
			b.Logf("Allocation size: %d bytes", size)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ptr := Malloc(size)
				Free(ptr)
			}
		})
	}
}

func BenchmarkMallocWithoutFree(b *testing.B) {
	b.ReportAllocs()
	var size Size = 64
	b.Logf("Allocation size: %d bytes (memory leak test)", size)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Malloc(size)
	}
}

func BenchmarkVoid(b *testing.B) {
	b.ReportAllocs()
	var size Size = 1024
	ptr := Malloc(size)
	defer Free(ptr)

	b.Logf("Zeroing %d bytes", size)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Void(ptr, size)
	}
}

func BenchmarkAllocateGeneric(b *testing.B) {
	b.ReportAllocs()

	type SmallStruct struct {
		a int
		b int
	}

	type MediumStruct struct {
		data [64]byte
	}

	type LargeStruct struct {
		data [1024]byte
	}

	b.Run("SmallStruct", func(b *testing.B) {
		size := SizeOf[SmallStruct]()
		b.Logf("SmallStruct size: %d bytes", size)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ptr := Allocate[SmallStruct]()
			DeAllocate(ptr)
		}
	})

	b.Run("MediumStruct", func(b *testing.B) {
		size := SizeOf[MediumStruct]()
		b.Logf("MediumStruct size: %d bytes", size)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ptr := Allocate[MediumStruct]()
			DeAllocate(ptr)
		}
	})

	b.Run("LargeStruct", func(b *testing.B) {
		size := SizeOf[LargeStruct]()
		b.Logf("LargeStruct size: %d bytes", size)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ptr := Allocate[LargeStruct]()
			DeAllocate(ptr)
		}
	})
}

func BenchmarkAllocateBlank(b *testing.B) {
	b.ReportAllocs()

	type DataStruct struct {
		id      int
		name    [32]byte
		payload [128]byte
	}

	size := SizeOf[DataStruct]()
	b.Logf("Allocating and zeroing %d bytes", size)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ptr := AllocateBlank[DataStruct]()
		DeAllocate(ptr)
	}
}

func BenchmarkSizeOf(b *testing.B) {
	b.ReportAllocs()

	type TestStruct struct {
		a int
		b uint
		c float64
		d [100]byte
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = SizeOf[TestStruct]()
	}
}

func BenchmarkAllocateVsAllocateBlank(b *testing.B) {
	b.ReportAllocs()

	type DataBlock struct {
		data [256]byte
	}

	size := SizeOf[DataBlock]()
	b.Logf("Comparing allocation methods for %d bytes", size)

	b.Run("Allocate", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ptr := Allocate[DataBlock]()
			DeAllocate(ptr)
		}
	})

	b.Run("AllocateBlank", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ptr := AllocateBlank[DataBlock]()
			DeAllocate(ptr)
		}
	})
}

func BenchmarkMultipleAllocations(b *testing.B) {
	b.ReportAllocs()

	type Node struct {
		value int
		next  *Node
	}

	nodeSize := SizeOf[Node]()
	b.Logf("Node size: %d bytes, allocating chains of varying lengths", nodeSize)

	allocations := []int{10, 100, 1000}

	for _, count := range allocations {
		b.Run(string(rune('0'+count/10)), func(b *testing.B) {
			b.Logf("Allocating %d nodes (%d total bytes)", count, count*int(nodeSize))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ptrs := make([]*Node, count)
				for j := 0; j < count; j++ {
					ptrs[j] = Allocate[Node]()
				}
				for j := 0; j < count; j++ {
					DeAllocate(ptrs[j])
				}
			}
		})
	}
}
