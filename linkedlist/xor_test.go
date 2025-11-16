package linkedlist

import (
	"runtime"
	"testing"
	"time"
)

func TestNewXORList(t *testing.T) {
	list := NewXORList[int]()
	defer list.Free()

	if list == nil {
		t.Fatal("NewXORList returned nil")
	}
	if !list.IsEmpty() {
		t.Errorf("New list should be empty")
	}
	if list.Size() != 0 {
		t.Errorf("New list size should be 0, got %d", list.Size())
	}
}

func TestXORList_Push(t *testing.T) {
	list := NewXORList[int]()

	list.Push(1)
	if list.Size() != 1 {
		t.Errorf("Expected size 1, got %d", list.Size())
	}
	if list.IsEmpty() {
		t.Error("List should not be empty after push")
	}

	list.Push(2)
	list.Push(3)
	if list.Size() != 3 {
		t.Errorf("Expected size 3, got %d", list.Size())
	}
}

func TestXORList_PushFront(t *testing.T) {
	list := NewXORList[int]()
	defer list.Free()

	// Test single PushFront
	list.PushFront(1)
	if list.Size() != 1 {
		t.Errorf("Expected size 1, got %d", list.Size())
	}
	if list.First() != 1 {
		t.Errorf("Expected first to be 1, got %d", list.First())
	}
	if list.Last() != 1 {
		t.Errorf("Expected last to be 1, got %d", list.Last())
	}

	// Test multiple PushFront calls - should insert at the beginning
	list.PushFront(2)
	list.PushFront(3)
	if list.Size() != 3 {
		t.Errorf("Expected size 3, got %d", list.Size())
	}

	// Elements should be in reverse order: 3, 2, 1
	if list.First() != 3 {
		t.Errorf("Expected first to be 3, got %d", list.First())
	}
	if list.At(1) != 2 {
		t.Errorf("Expected At(1) to be 2, got %d", list.At(1))
	}
	if list.Last() != 1 {
		t.Errorf("Expected last to be 1, got %d", list.Last())
	}
}

func TestXORList_PushFrontWithPopRight(t *testing.T) {
	list := NewXORList[int]()
	defer list.Free()

	// PushFront adds to the beginning, PopRight removes from the end
	list.PushFront(1)
	list.PushFront(2)
	list.PushFront(3)
	// List is now: 3, 2, 1

	val := list.PopRight()
	if val != 1 {
		t.Errorf("PopRight expected 1, got %d", val)
	}
	if list.Size() != 2 {
		t.Errorf("Expected size 2, got %d", list.Size())
	}

	val = list.PopRight()
	if val != 2 {
		t.Errorf("PopRight expected 2, got %d", val)
	}
	if list.Size() != 1 {
		t.Errorf("Expected size 1, got %d", list.Size())
	}

	val = list.PopRight()
	if val != 3 {
		t.Errorf("PopRight expected 3, got %d", val)
	}
	if !list.IsEmpty() {
		t.Error("List should be empty after all PopRight operations")
	}
}

func TestXORList_PushFrontWithPopLeft(t *testing.T) {
	list := NewXORList[int]()
	defer list.Free()

	// PushFront adds to the beginning, PopLeft removes from the beginning
	list.PushFront(1)
	list.PushFront(2)
	list.PushFront(3)
	// List is now: 3, 2, 1

	val := list.PopLeft()
	if val != 3 {
		t.Errorf("PopLeft expected 3, got %d", val)
	}
	if list.Size() != 2 {
		t.Errorf("Expected size 2, got %d", list.Size())
	}

	val = list.PopLeft()
	if val != 2 {
		t.Errorf("PopLeft expected 2, got %d", val)
	}
	if list.Size() != 1 {
		t.Errorf("Expected size 1, got %d", list.Size())
	}

	val = list.PopLeft()
	if val != 1 {
		t.Errorf("PopLeft expected 1, got %d", val)
	}
	if !list.IsEmpty() {
		t.Error("List should be empty after all PopLeft operations")
	}
}

func TestXORList_PushFrontAndPushMixed(t *testing.T) {
	list := NewXORList[int]()
	defer list.Free()

	// Mix PushFront and Push operations
	list.Push(1)      // [1]
	list.PushFront(2) // [2, 1]
	list.Push(3)      // [2, 1, 3]
	list.PushFront(4) // [4, 2, 1, 3]

	if list.Size() != 4 {
		t.Errorf("Expected size 4, got %d", list.Size())
	}

	expected := []int{4, 2, 1, 3}
	for i, exp := range expected {
		if got := list.At(i); got != exp {
			t.Errorf("At(%d) = %d, want %d", i, got, exp)
		}
	}

	// Test PopRight on mixed list
	val := list.PopRight()
	if val != 3 {
		t.Errorf("PopRight expected 3, got %d", val)
	}

	// Test PopLeft on mixed list
	val = list.PopLeft()
	if val != 4 {
		t.Errorf("PopLeft expected 4, got %d", val)
	}

	if list.Size() != 2 {
		t.Errorf("Expected size 2, got %d", list.Size())
	}
}

func TestXORList_PushAll(t *testing.T) {
	list := NewXORList[int]()
	defer list.Free()

	list.PushAll(1, 2, 3, 4, 5)
	if list.Size() != 5 {
		t.Errorf("Expected size 5, got %d", list.Size())
	}

	if list.At(0) != 1 || list.At(4) != 5 {
		t.Error("Elements not in expected order")
	}
}

func TestXORList_FirstLast(t *testing.T) {
	list := NewXORList[int]()
	defer list.Free()

	if list.First() != 0 {
		t.Error("First on empty list should return zero value")
	}
	if list.Last() != 0 {
		t.Error("Last on empty list should return zero value")
	}

	list.Push(42)
	if list.First() != 42 {
		t.Errorf("Expected first to be 42, got %d", list.First())
	}
	if list.Last() != 42 {
		t.Errorf("Expected last to be 42, got %d", list.Last())
	}

	list.Push(100)
	list.Push(200)
	if list.First() != 42 {
		t.Errorf("Expected first to be 42, got %d", list.First())
	}
	if list.Last() != 200 {
		t.Errorf("Expected last to be 200, got %d", list.Last())
	}
}

func TestXORList_At(t *testing.T) {
	list := NewXORList[int]()
	defer list.Free()

	list.PushAll(10, 20, 30, 40, 50)

	tests := []struct {
		index    int
		expected int
	}{
		{0, 10},
		{1, 20},
		{2, 30},
		{3, 40},
		{4, 50},
	}

	for _, tt := range tests {
		if got := list.At(tt.index); got != tt.expected {
			t.Errorf("At(%d) = %d, want %d", tt.index, got, tt.expected)
		}
	}
}

func TestXORList_Get(t *testing.T) {
	list := NewXORList[string]()
	defer list.Free()

	list.PushAll("a", "b", "c")

	if list.Get(0) != "a" {
		t.Errorf("Get(0) expected 'a', got '%s'", list.Get(0))
	}
	if list.Get(2) != "c" {
		t.Errorf("Get(2) expected 'c', got '%s'", list.Get(2))
	}
}

func TestXORList_PopLeft(t *testing.T) {
	list := NewXORList[int]()
	defer list.Free()

	if list.PopLeft() != 0 {
		t.Error("PopLeft on empty list should return zero value")
	}

	list.PushAll(1, 2, 3, 4, 5)

	val := list.PopLeft()
	if val != 1 {
		t.Errorf("PopLeft expected 1, got %d", val)
	}
	if list.Size() != 4 {
		t.Errorf("Size should be 4 after PopLeft, got %d", list.Size())
	}
	if list.First() != 2 {
		t.Errorf("First should be 2 after PopLeft, got %d", list.First())
	}
}

func TestXORList_PopRight(t *testing.T) {
	list := NewXORList[int]()
	defer list.Free()

	if list.PopRight() != 0 {
		t.Error("PopRight on empty list should return zero value")
	}

	list.PushAll(1, 2, 3, 4, 5)

	val := list.PopRight()
	if val != 5 {
		t.Errorf("PopRight expected 5, got %d", val)
	}
	if list.Size() != 4 {
		t.Errorf("Size should be 4 after PopRight, got %d", list.Size())
	}
	if list.Last() != 4 {
		t.Errorf("Last should be 4 after PopRight, got %d", list.Last())
	}
}

func TestXORList_Delete(t *testing.T) {
	list := NewXORList[int]()
	defer list.Free()

	list.PushAll(1, 2, 3, 4, 5)

	list.Delete(2)
	if list.Size() != 4 {
		t.Errorf("Expected size 4, got %d", list.Size())
	}

	expected := []int{1, 2, 4, 5}
	for i, exp := range expected {
		if got := list.At(i); got != exp {
			t.Errorf("At(%d) = %d, want %d", i, got, exp)
		}
	}
}

func TestXORList_DeleteBy(t *testing.T) {
	list := NewXORList[int]()
	defer list.Free()

	list.PushAll(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	list.DeleteBy(func(val int) bool {
		return val%2 == 0
	})

	if list.Size() != 5 {
		t.Errorf("Expected size 5 after deleting evens, got %d", list.Size())
	}

	list.ForEach(func(index int, data int) bool {
		if data%2 == 0 {
			t.Errorf("Found even number %d after DeleteBy", data)
		}
		return true
	})
}

func TestXORList_DeleteAll(t *testing.T) {
	list := NewXORList[int]()
	defer list.Free()

	list.PushAll(1, 2, 3, 4, 5)

	list.DeleteAll()

	if !list.IsEmpty() {
		t.Error("List should be empty after DeleteAll")
	}
	if list.Size() != 0 {
		t.Errorf("Size should be 0 after DeleteAll, got %d", list.Size())
	}
}

func TestXORList_ForEach(t *testing.T) {
	list := NewXORList[int]()
	defer list.Free()

	list.PushAll(1, 2, 3, 4, 5)

	sum := 0
	list.ForEach(func(index int, data int) bool {
		sum += data
		return true
	})

	if sum != 15 {
		t.Errorf("Expected sum 15, got %d", sum)
	}

	count := 0
	list.ForEach(func(index int, data int) bool {
		count++
		return count < 3
	})

	if count != 3 {
		t.Errorf("Expected count 3 (early termination), got %d", count)
	}
}

func TestXORList_Find(t *testing.T) {
	list := NewXORList[int]()
	defer list.Free()

	list.PushAll(1, 2, 3, 4, 5)

	val, found := list.Find(func(v int) bool {
		return v == 3
	})
	if !found {
		t.Error("Expected to find value 3")
	}
	if val != 3 {
		t.Errorf("Expected value 3, got %d", val)
	}

	_, found = list.Find(func(v int) bool {
		return v == 99
	})
	if found {
		t.Error("Should not find value 99")
	}
}

func TestXORList_Some(t *testing.T) {
	list := NewXORList[int]()
	defer list.Free()

	list.PushAll(1, 2, 3, 4, 5)

	if !list.Some(func(v int) bool { return v%2 == 0 }) {
		t.Error("Some should return true for even numbers")
	}

	if list.Some(func(v int) bool { return v > 10 }) {
		t.Error("Some should return false for numbers > 10")
	}
}

func TestXORList_Filter(t *testing.T) {
	list := NewXORList[int]()
	defer list.Free()

	list.PushAll(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	filtered := list.Filter(func(v int) bool {
		return v%2 == 0
	})

	if filtered.Size() != 5 {
		t.Errorf("Expected filtered size 5, got %d", filtered.Size())
	}

	filtered.ForEach(func(index int, data int) bool {
		if data%2 != 0 {
			t.Errorf("Found odd number %d in filtered list", data)
		}
		return true
	})
}

func TestXORList_Join(t *testing.T) {
	list1 := NewXORList[int]()
	defer list1.Free()

	list1.PushAll(1, 2, 3)

	list2 := NewXORList[int]()
	defer list2.Free()

	list2.PushAll(4, 5, 6)

	list1.Join(list2)

	if list1.Size() != 6 {
		t.Errorf("Expected size 6 after join, got %d", list1.Size())
	}

	expected := []int{1, 2, 3, 4, 5, 6}
	for i, exp := range expected {
		if got := list1.At(i); got != exp {
			t.Errorf("At(%d) = %d, want %d", i, got, exp)
		}
	}
}

func TestXORList_Merge(t *testing.T) {
	list1 := NewXORList[int]()
	defer list1.Free()

	list1.PushAll(1, 2, 3)

	list2 := NewXORList[int]()
	defer list2.Free()

	list2.PushAll(4, 5, 6)

	merged := list1.Merge(list2)
	defer merged.DeleteAll()

	if merged.Size() != 6 {
		t.Errorf("Expected merged size 6, got %d", merged.Size())
	}

	if list1.Size() != 3 {
		t.Errorf("list1 should still have size 3, got %d", list1.Size())
	}
	if list2.Size() != 3 {
		t.Errorf("list2 should still have size 3, got %d", list2.Size())
	}
}

func TestXORList_Shrink(t *testing.T) {
	list := NewXORList[int]()
	defer list.Free()

	list.PushAll(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	list.Shrink(5)
	if list.Size() != 5 {
		t.Errorf("Expected size 5 after shrink, got %d", list.Size())
	}

	for i := 0; i < 5; i++ {
		if list.At(i) != i+1 {
			t.Errorf("At(%d) = %d, want %d", i, list.At(i), i+1)
		}
	}

	list.Shrink(0)
	if !list.IsEmpty() {
		t.Error("List should be empty after shrink to 0")
	}
}

func TestXORList_EdgeCases(t *testing.T) {
	list := NewXORList[int]()
	defer list.Free()

	list.Push(42)
	if list.First() != 42 || list.Last() != 42 {
		t.Error("Single element list First/Last mismatch")
	}

	val := list.PopLeft()
	if val != 42 {
		t.Errorf("PopLeft expected 42, got %d", val)
	}
	if !list.IsEmpty() {
		t.Error("List should be empty after popping single element")
	}

	list.Push(1)
	list.Push(2)
	if list.Size() != 2 {
		t.Errorf("Expected size 2, got %d", list.Size())
	}

	list.PopRight()
	if list.Size() != 1 || list.First() != 1 {
		t.Error("List state incorrect after PopRight on 2-element list")
	}
}

func TestXORList_StringType(t *testing.T) {
	list := NewXORList[string]()
	defer list.Free()

	list.PushAll("hello", "world", "test")

	if list.Size() != 3 {
		t.Errorf("Expected size 3, got %d", list.Size())
	}

	if list.First() != "hello" {
		t.Errorf("Expected first to be 'hello', got '%s'", list.First())
	}

	if list.Last() != "test" {
		t.Errorf("Expected last to be 'test', got '%s'", list.Last())
	}

	val, found := list.Find(func(s string) bool {
		return s == "world"
	})
	if !found || val != "world" {
		t.Error("Failed to find 'world' in list")
	}
}

func TestXORList_GarbageCollection(t *testing.T) {
	func() {
		list1 := NewXORList[int]()
		defer list1.Free()

		list1.PushAll(1, 2, 3, 4, 5)

		list2 := NewXORList[string]()
		defer list2.Free()

		list2.PushAll("a", "b", "c")

	}()

	runtime.GC()
	runtime.Gosched()
	time.Sleep(10 * time.Millisecond)

	runtime.GC()
	runtime.Gosched()
	time.Sleep(10 * time.Millisecond)

	list3 := NewXORList[int]()
	defer list3.Free()

	list3.PushAll(10, 20, 30, 40, 50)

	if list3.Size() != 5 {
		t.Errorf("Expected size 5, got %d", list3.Size())
	}
	if list3.First() != 10 {
		t.Errorf("Expected first to be 10, got %d", list3.First())
	}
	if list3.Last() != 50 {
		t.Errorf("Expected last to be 50, got %d", list3.Last())
	}

	list3.DeleteAll()

	runtime.GC()
	runtime.Gosched()
	time.Sleep(5 * time.Millisecond)
}

func TestXORList_NoFinalizerInterference(t *testing.T) {
	for i := 0; i < 10; i++ {
		list := NewXORList[int]()

		list.PushAll(1, 2, 3, 4, 5)

		if list.Size() != 5 {
			t.Errorf("Iteration %d: expected size 5, got %d", i, list.Size())
		}

		list.Free()
	}

	runtime.GC()
	runtime.Gosched()
}

func TestXORList_ConcurrentGC(t *testing.T) {
	list := NewXORList[int]()
	defer list.Free()

	for i := 0; i < 100; i++ {
		list.Push(i)

		if i%10 == 0 {
			runtime.GC()
			runtime.Gosched()
		}
	}

	if list.Size() != 100 {
		t.Errorf("Expected size 100, got %d", list.Size())
	}

	for i := 0; i < 50; i++ {
		list.PopLeft()

		if i%10 == 0 {
			runtime.GC()
			runtime.Gosched()
		}
	}

	if list.Size() != 50 {
		t.Errorf("Expected size 50 after pops, got %d", list.Size())
	}
}

// Benchmark tests with allocation tracing

func BenchmarkXORList_Push(b *testing.B) {
	b.ReportAllocs()

	type benchStruct struct {
		id   int
		data [32]byte
	}

	b.Run("int", func(b *testing.B) {
		list := NewXORList[int]()
		defer list.Free()

		b.Logf("Pushing %d integers (size: %d bytes each)", b.N, 8)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			list.Push(i)
		}
	})

	b.Run("string", func(b *testing.B) {
		list := NewXORList[string]()
		defer list.Free()

		b.Logf("Pushing %d strings", b.N)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			list.Push("test string data")
		}
	})

	b.Run("struct", func(b *testing.B) {
		list := NewXORList[benchStruct]()
		defer list.Free()

		b.Logf("Pushing %d structs (size: ~40 bytes each)", b.N)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			list.Push(benchStruct{id: i})
		}
	})
}

func BenchmarkXORList_PushAll(b *testing.B) {
	b.ReportAllocs()

	sizes := []int{10, 100, 1000}

	for _, size := range sizes {
		b.Run(string(rune('0'+size/10)), func(b *testing.B) {
			data := make([]int, size)
			for i := range data {
				data[i] = i
			}

			b.Logf("PushAll %d elements, total memory: ~%d bytes", size, size*24)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				list := NewXORList[int]()
				list.PushAll(data...)
				list.Free()
			}
		})
	}
}

func BenchmarkXORList_PopLeft(b *testing.B) {
	b.ReportAllocs()

	sizes := []int{100, 1000, 10000}

	for _, size := range sizes {
		b.Run(string(rune('0'+size/100)), func(b *testing.B) {
			b.StopTimer()
			lists := make([]*XORLinkedList[int, int], b.N)
			for i := 0; i < b.N; i++ {
				lists[i] = NewXORList[int]()
				for j := 0; j < size; j++ {
					lists[i].Push(j)
				}
			}

			b.Logf("PopLeft from lists of %d elements", size)
			b.StartTimer()
			for i := 0; i < b.N; i++ {
				for j := 0; j < size; j++ {
					lists[i].PopLeft()
				}
				lists[i].Free()
			}
		})
	}
}

func BenchmarkXORList_PopRight(b *testing.B) {
	b.ReportAllocs()

	sizes := []int{100, 1000, 10000}

	for _, size := range sizes {
		b.Run(string(rune('0'+size/100)), func(b *testing.B) {
			b.StopTimer()
			lists := make([]*XORLinkedList[int, int], b.N)
			for i := 0; i < b.N; i++ {
				lists[i] = NewXORList[int]()
				for j := 0; j < size; j++ {
					lists[i].Push(j)
				}
			}

			b.Logf("PopRight from lists of %d elements", size)
			b.StartTimer()
			for i := 0; i < b.N; i++ {
				for j := 0; j < size; j++ {
					lists[i].PopRight()
				}
				lists[i].Free()
			}
		})
	}
}

func BenchmarkXORList_At(b *testing.B) {
	b.ReportAllocs()

	list := NewXORList[int]()
	defer list.Free()

	size := 1000
	for i := 0; i < size; i++ {
		list.Push(i)
	}

	b.Logf("Random access in list of %d elements", size)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = list.At(i % size)
	}
}

func BenchmarkXORList_Delete(b *testing.B) {
	b.ReportAllocs()

	positions := []string{"front", "middle", "back"}

	for _, pos := range positions {
		b.Run(pos, func(b *testing.B) {
			b.StopTimer()
			size := 1000

			lists := make([]*XORLinkedList[int, int], b.N)
			indices := make([]int, b.N)

			for i := 0; i < b.N; i++ {
				lists[i] = NewXORList[int]()
				for j := 0; j < size; j++ {
					lists[i].Push(j)
				}

				switch pos {
				case "front":
					indices[i] = 0
				case "middle":
					indices[i] = size / 2
				case "back":
					indices[i] = size - 1
				}
			}
			b.Logf("Delete from %s of %d element list", pos, size)
			b.StartTimer()
			for i := 0; i < b.N; i++ {
				lists[i].Delete(indices[i])
				lists[i].Free()
			}
		})
	}
}

func BenchmarkXORList_DeleteBy(b *testing.B) {
	b.ReportAllocs()

	size := 1000
	b.Logf("DeleteBy (removing evens) from %d element list", size)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		list := NewXORList[int]()
		for j := 0; j < size; j++ {
			list.Push(j)
		}
		b.StartTimer()

		list.DeleteBy(func(v int) bool {
			return v%2 == 0
		})

		b.StopTimer()
		list.Free()
		b.StartTimer()
	}
}

func BenchmarkXORList_ForEach(b *testing.B) {
	b.ReportAllocs()

	sizes := []int{100, 1000, 10000}

	for _, size := range sizes {
		b.Run(string(rune('0'+size/100)), func(b *testing.B) {
			list := NewXORList[int]()
			defer list.Free()

			for i := 0; i < size; i++ {
				list.Push(i)
			}

			b.Logf("ForEach over %d elements", size)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				sum := 0
				list.ForEach(func(index int, data int) bool {
					sum += data
					return true
				})
			}
		})
	}
}

func BenchmarkXORList_Find(b *testing.B) {
	b.ReportAllocs()

	list := NewXORList[int]()
	defer list.Free()

	size := 1000
	for i := 0; i < size; i++ {
		list.Push(i)
	}

	b.Logf("Find in list of %d elements", size)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list.Find(func(v int) bool {
			return v == size/2
		})
	}
}

func BenchmarkXORList_Filter(b *testing.B) {
	b.ReportAllocs()

	size := 1000
	b.Logf("Filter (selecting evens) from %d element list", size)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		list := NewXORList[int]()
		for j := 0; j < size; j++ {
			list.Push(j)
		}
		b.StartTimer()

		filtered := list.Filter(func(v int) bool {
			return v%2 == 0
		})

		b.StopTimer()
		list.Free()
		filtered.DeleteAll()
		b.StartTimer()
	}
}

func BenchmarkXORList_Join(b *testing.B) {
	b.ReportAllocs()

	sizes := []int{100, 500, 1000}

	for _, size := range sizes {
		b.Run(string(rune('0'+size/100)), func(b *testing.B) {
			b.Logf("Joining two lists of %d elements each", size)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				list1 := NewXORList[int]()
				list2 := NewXORList[int]()

				for j := 0; j < size; j++ {
					list1.Push(j)
					list2.Push(j + size)
				}
				b.StartTimer()

				list1.Join(list2)

				b.StopTimer()
				list1.Free()
				list2.Free()
				b.StartTimer()
			}
		})
	}
}

func BenchmarkXORList_Merge(b *testing.B) {
	b.ReportAllocs()

	size := 500
	b.Logf("Merging two lists of %d elements each", size)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		list1 := NewXORList[int]()
		list2 := NewXORList[int]()

		for j := 0; j < size; j++ {
			list1.Push(j)
			list2.Push(j + size)
		}
		b.StartTimer()

		merged := list1.Merge(list2)

		b.StopTimer()
		list1.Free()
		list2.Free()
		merged.DeleteAll()
		b.StartTimer()
	}
}

func BenchmarkXORList_Shrink(b *testing.B) {
	b.ReportAllocs()

	size := 1000
	shrinkTo := 100

	b.Logf("Shrinking list from %d to %d elements", size, shrinkTo)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		list := NewXORList[int]()
		for j := 0; j < size; j++ {
			list.Push(j)
		}
		b.StartTimer()

		list.Shrink(shrinkTo)

		b.StopTimer()
		list.Free()
		b.StartTimer()
	}
}

func BenchmarkXORList_MemoryFootprint(b *testing.B) {
	b.ReportAllocs()

	type LargeStruct struct {
		id      int
		name    [64]byte
		payload [256]byte
	}

	sizes := []int{10, 100, 1000}

	for _, size := range sizes {
		b.Run(string(rune('0'+size/10)), func(b *testing.B) {
			// XOR node overhead is minimal (just the xor link + data)
			nodeOverhead := 8 // uintptr for xor link
			dataSize := 328   // LargeStruct size
			totalPerNode := nodeOverhead + dataSize

			b.Logf("Creating list of %d LargeStruct elements (total memory: ~%d bytes)",
				size, size*totalPerNode)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				list := NewXORList[LargeStruct]()
				for j := 0; j < size; j++ {
					list.Push(LargeStruct{id: j})
				}
				list.Free()
			}
		})
	}
}

func BenchmarkXORList_CreateAndDestroy(b *testing.B) {
	b.ReportAllocs()

	sizes := []int{10, 100, 1000}

	for _, size := range sizes {
		b.Run(string(rune('0'+size/10)), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				list := NewXORList[int]()
				for j := 0; j < size; j++ {
					list.Push(j)
				}
				list.Free()
			}
		})
	}
}
