package linkedlist

import (
	"fmt"
	"testing"
)

// ============================================================================
// COMPREHENSIVE BENCHMARK SUITE FOR LINKED LIST
// ============================================================================

// ----------------------------------------------------------------------------
// Benchmarks: Push Operations
// ----------------------------------------------------------------------------

func BenchmarkLinkedList_Push_Small(b *testing.B) {
	benchmarkPush(b, 10)
}

func BenchmarkLinkedList_Push_Medium(b *testing.B) {
	benchmarkPush(b, 100)
}

func BenchmarkLinkedList_Push_Large(b *testing.B) {
	benchmarkPush(b, 1000)
}

func BenchmarkLinkedList_Push_ExtraLarge(b *testing.B) {
	benchmarkPush(b, 10000)
}

func benchmarkPush(b *testing.B, size int) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list := NewLinkedList[int]()
		for j := 0; j < size; j++ {
			list.Push(j)
		}
	}
}

// ----------------------------------------------------------------------------
// Benchmarks: PushFront Operations
// ----------------------------------------------------------------------------

func BenchmarkLinkedList_PushFront_Small(b *testing.B) {
	benchmarkPushFront(b, 10)
}

func BenchmarkLinkedList_PushFront_Medium(b *testing.B) {
	benchmarkPushFront(b, 100)
}

func BenchmarkLinkedList_PushFront_Large(b *testing.B) {
	benchmarkPushFront(b, 1000)
}

func BenchmarkLinkedList_PushFront_ExtraLarge(b *testing.B) {
	benchmarkPushFront(b, 10000)
}

func benchmarkPushFront(b *testing.B, size int) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list := NewLinkedList[int]()
		for j := 0; j < size; j++ {
			list.PushFront(j)
		}
	}
}

// ----------------------------------------------------------------------------
// Benchmarks: At/Get Operations (Random Access)
// ----------------------------------------------------------------------------

func BenchmarkLinkedList_At_First(b *testing.B) {
	list := NewLinkedList[int]()
	for i := 0; i < 1000; i++ {
		list.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = list.At(0)
	}
}

func BenchmarkLinkedList_At_Middle(b *testing.B) {
	list := NewLinkedList[int]()
	for i := 0; i < 1000; i++ {
		list.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = list.At(500)
	}
}

func BenchmarkLinkedList_At_Last(b *testing.B) {
	list := NewLinkedList[int]()
	for i := 0; i < 1000; i++ {
		list.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = list.At(999)
	}
}

func BenchmarkLinkedList_At_NegativeIndex(b *testing.B) {
	list := NewLinkedList[int]()
	for i := 0; i < 1000; i++ {
		list.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = list.At(-1)
	}
}

func BenchmarkLinkedList_At_Sequential(b *testing.B) {
	list := NewLinkedList[int]()
	for i := 0; i < 1000; i++ {
		list.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1000; j++ {
			_ = list.At(j)
		}
	}
}

// ----------------------------------------------------------------------------
// Benchmarks: PopLeft Operations
// ----------------------------------------------------------------------------

func BenchmarkLinkedList_PopLeft_Small(b *testing.B) {
	benchmarkPopLeft(b, 10)
}

func BenchmarkLinkedList_PopLeft_Medium(b *testing.B) {
	benchmarkPopLeft(b, 100)
}

func BenchmarkLinkedList_PopLeft_Large(b *testing.B) {
	benchmarkPopLeft(b, 1000)
}

func benchmarkPopLeft(b *testing.B, size int) {
	b.ResetTimer()
	b.Logf("b.N: %v", b.N)
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		list := NewLinkedList[int]()
		for j := 0; j < size; j++ {
			list.Push(j)
		}
		b.StartTimer()

		for j := 0; j < size; j++ {
			list.PopLeft()
		}
	}
}

// ----------------------------------------------------------------------------
// Benchmarks: PopRight Operations
// ----------------------------------------------------------------------------

func BenchmarkLinkedList_PopRight_Small(b *testing.B) {
	benchmarkPopRight(b, 10)
}

func BenchmarkLinkedList_PopRight_Medium(b *testing.B) {
	benchmarkPopRight(b, 100)
}

func BenchmarkLinkedList_PopRight_Large(b *testing.B) {
	benchmarkPopRight(b, 1000)
}

func benchmarkPopRight(b *testing.B, size int) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		list := NewLinkedList[int]()
		for j := 0; j < size; j++ {
			list.Push(j)
		}
		b.StartTimer()

		for j := 0; j < size; j++ {
			list.PopRight()
		}
	}
}

// ----------------------------------------------------------------------------
// Benchmarks: Pop (by index) Operations
// ----------------------------------------------------------------------------

func BenchmarkLinkedList_Pop_First(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		list := NewLinkedList[int]()
		for j := 0; j < 100; j++ {
			list.Push(j)
		}
		b.StartTimer()
		list.Pop(0)
	}
}

func BenchmarkLinkedList_Pop_Middle(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		list := NewLinkedList[int]()
		for j := 0; j < 100; j++ {
			list.Push(j)
		}
		b.StartTimer()
		list.Pop(50)
	}
}

func BenchmarkLinkedList_Pop_Last(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		list := NewLinkedList[int]()
		for j := 0; j < 100; j++ {
			list.Push(j)
		}
		b.StartTimer()
		list.Pop(99)
	}
}

// ----------------------------------------------------------------------------
// Benchmarks: Delete Operations
// ----------------------------------------------------------------------------

func BenchmarkLinkedList_Delete_First(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		list := NewLinkedList[int]()
		for j := 0; j < 100; j++ {
			list.Push(j)
		}
		b.StartTimer()
		list.Delete(0)
	}
}

func BenchmarkLinkedList_Delete_Middle(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		list := NewLinkedList[int]()
		for j := 0; j < 100; j++ {
			list.Push(j)
		}
		b.StartTimer()
		list.Delete(50)
	}
}

func BenchmarkLinkedList_Delete_Last(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		list := NewLinkedList[int]()
		for j := 0; j < 100; j++ {
			list.Push(j)
		}
		b.StartTimer()
		list.Delete(99)
	}
}

func BenchmarkLinkedList_DeleteBy(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		list := NewLinkedList[int]()
		for j := 0; j < 1000; j++ {
			list.Push(j)
		}
		b.StartTimer()
		list.DeleteBy(func(v int) bool {
			return v%2 == 0
		})
	}
}

func BenchmarkLinkedList_DeleteAll(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		list := NewLinkedList[int]()
		for j := 0; j < 1000; j++ {
			list.Push(j)
		}
		b.StartTimer()
		list.DeleteAll()
	}
}

// ----------------------------------------------------------------------------
// Benchmarks: Search Operations
// ----------------------------------------------------------------------------

func BenchmarkLinkedList_IndexOf_Found_First(b *testing.B) {
	list := NewLinkedList[int]()
	for i := 0; i < 1000; i++ {
		list.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list.IndexOf(func(v int) bool {
			return v == 0
		})
	}
}

func BenchmarkLinkedList_IndexOf_Found_Middle(b *testing.B) {
	list := NewLinkedList[int]()
	for i := 0; i < 1000; i++ {
		list.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list.IndexOf(func(v int) bool {
			return v == 500
		})
	}
}

func BenchmarkLinkedList_IndexOf_Found_Last(b *testing.B) {
	list := NewLinkedList[int]()
	for i := 0; i < 1000; i++ {
		list.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list.IndexOf(func(v int) bool {
			return v == 999
		})
	}
}

func BenchmarkLinkedList_IndexOf_NotFound(b *testing.B) {
	list := NewLinkedList[int]()
	for i := 0; i < 1000; i++ {
		list.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list.IndexOf(func(v int) bool {
			return v == 9999
		})
	}
}

func BenchmarkLinkedList_Find_Found_Early(b *testing.B) {
	list := NewLinkedList[int]()
	for i := 0; i < 1000; i++ {
		list.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list.Find(func(v int) bool {
			return v == 10
		})
	}
}

func BenchmarkLinkedList_Find_NotFound(b *testing.B) {
	list := NewLinkedList[int]()
	for i := 0; i < 1000; i++ {
		list.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list.Find(func(v int) bool {
			return v == 9999
		})
	}
}

func BenchmarkLinkedList_Some_Found_Early(b *testing.B) {
	list := NewLinkedList[int]()
	for i := 0; i < 1000; i++ {
		list.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list.Some(func(v int) bool {
			return v == 10
		})
	}
}

func BenchmarkLinkedList_Some_NotFound(b *testing.B) {
	list := NewLinkedList[int]()
	for i := 0; i < 1000; i++ {
		list.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list.Some(func(v int) bool {
			return v > 10000
		})
	}
}

// ----------------------------------------------------------------------------
// Benchmarks: Filter Operations
// ----------------------------------------------------------------------------

func BenchmarkLinkedList_Filter_Half(b *testing.B) {
	list := NewLinkedList[int]()
	for i := 0; i < 1000; i++ {
		list.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list.Filter(func(v int) bool {
			return v%2 == 0
		})
	}
}

func BenchmarkLinkedList_Filter_Few(b *testing.B) {
	list := NewLinkedList[int]()
	for i := 0; i < 1000; i++ {
		list.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list.Filter(func(v int) bool {
			return v < 10
		})
	}
}

func BenchmarkLinkedList_Filter_All(b *testing.B) {
	list := NewLinkedList[int]()
	for i := 0; i < 1000; i++ {
		list.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list.Filter(func(v int) bool {
			return true
		})
	}
}

func BenchmarkLinkedList_Filter_None(b *testing.B) {
	list := NewLinkedList[int]()
	for i := 0; i < 1000; i++ {
		list.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list.Filter(func(v int) bool {
			return false
		})
	}
}

// ----------------------------------------------------------------------------
// Benchmarks: ForEach Operations
// ----------------------------------------------------------------------------

func BenchmarkLinkedList_ForEach_Complete(b *testing.B) {
	list := NewLinkedList[int]()
	for i := 0; i < 1000; i++ {
		list.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum := 0
		list.ForEach(func(idx int, val int) bool {
			sum += val
			return true
		})
	}
}

func BenchmarkLinkedList_ForEach_EarlyExit(b *testing.B) {
	list := NewLinkedList[int]()
	for i := 0; i < 1000; i++ {
		list.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list.ForEach(func(idx int, val int) bool {
			return val < 10
		})
	}
}

// ----------------------------------------------------------------------------
// Benchmarks: Swap Operations
// ----------------------------------------------------------------------------

func BenchmarkLinkedList_Swap_Adjacent(b *testing.B) {
	list := NewLinkedList[int]()
	for i := 0; i < 100; i++ {
		list.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list.Swap(10, 11)
	}
}

func BenchmarkLinkedList_Swap_Distant(b *testing.B) {
	list := NewLinkedList[int]()
	for i := 0; i < 100; i++ {
		list.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list.Swap(10, 90)
	}
}

func BenchmarkLinkedList_Swap_FirstAndLast(b *testing.B) {
	list := NewLinkedList[int]()
	for i := 0; i < 100; i++ {
		list.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list.Swap(0, 99)
	}
}

// ----------------------------------------------------------------------------
// Benchmarks: Move Operations
// ----------------------------------------------------------------------------

func BenchmarkLinkedList_Move_Forward(b *testing.B) {
	list := NewLinkedList[int]()
	for i := 0; i < 100; i++ {
		list.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list.Move(10, 90)
	}
}

func BenchmarkLinkedList_Move_Backward(b *testing.B) {
	list := NewLinkedList[int]()
	for i := 0; i < 100; i++ {
		list.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list.Move(90, 10)
	}
}

func BenchmarkLinkedList_Move_Adjacent(b *testing.B) {
	list := NewLinkedList[int]()
	for i := 0; i < 100; i++ {
		list.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list.Move(10, 11)
	}
}

// ----------------------------------------------------------------------------
// Benchmarks: Join and Merge Operations
// ----------------------------------------------------------------------------

func BenchmarkLinkedList_Join_Small(b *testing.B) {
	benchmarkJoin(b, 10, 10)
}

func BenchmarkLinkedList_Join_Medium(b *testing.B) {
	benchmarkJoin(b, 100, 100)
}

func BenchmarkLinkedList_Join_Large(b *testing.B) {
	benchmarkJoin(b, 1000, 1000)
}

func benchmarkJoin(b *testing.B, size1, size2 int) {
	list2 := NewLinkedList[int]()
	for i := 0; i < size2; i++ {
		list2.Push(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		list1 := NewLinkedList[int]()
		for j := 0; j < size1; j++ {
			list1.Push(j)
		}
		b.StartTimer()
		list1.Join(list2)
	}
}

func BenchmarkLinkedList_Merge_Small(b *testing.B) {
	benchmarkMerge(b, 10, 10)
}

func BenchmarkLinkedList_Merge_Medium(b *testing.B) {
	benchmarkMerge(b, 100, 100)
}

func BenchmarkLinkedList_Merge_Large(b *testing.B) {
	benchmarkMerge(b, 1000, 1000)
}

func benchmarkMerge(b *testing.B, size1, size2 int) {
	list1 := NewLinkedList[int]()
	for i := 0; i < size1; i++ {
		list1.Push(i)
	}
	list2 := NewLinkedList[int]()
	for i := 0; i < size2; i++ {
		list2.Push(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list1.Merge(list2)
	}
}

// ----------------------------------------------------------------------------
// Benchmarks: Shrink Operations
// ----------------------------------------------------------------------------

func BenchmarkLinkedList_Shrink_Half(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		list := NewLinkedList[int]()
		for j := 0; j < 1000; j++ {
			list.Push(j)
		}
		b.StartTimer()
		list.Shrink(500)
	}
}

func BenchmarkLinkedList_Shrink_MostElements(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		list := NewLinkedList[int]()
		for j := 0; j < 1000; j++ {
			list.Push(j)
		}
		b.StartTimer()
		list.Shrink(10)
	}
}

func BenchmarkLinkedList_Shrink_ToZero(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		list := NewLinkedList[int]()
		for j := 0; j < 1000; j++ {
			list.Push(j)
		}
		b.StartTimer()
		list.Shrink(0)
	}
}

// ----------------------------------------------------------------------------
// Benchmarks: Mixed Operations
// ----------------------------------------------------------------------------

func BenchmarkLinkedList_Mixed_Realistic(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list := NewLinkedList[int]()

		// Add elements
		for j := 0; j < 100; j++ {
			list.Push(j)
		}

		// Mix of operations
		for j := 0; j < 50; j++ {
			_ = list.At(j)
			list.Push(j + 100)
			if j%5 == 0 {
				list.PopLeft()
			}
			if j%7 == 0 {
				list.Delete(j % list.Size())
			}
		}

		// Clean up
		list.DeleteAll()
	}
}

func BenchmarkLinkedList_Mixed_PushPopIntensive(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list := NewLinkedList[int]()

		for j := 0; j < 1000; j++ {
			list.Push(j)
			if j%3 == 0 && list.Size() > 0 {
				list.PopLeft()
			}
			if j%5 == 0 && list.Size() > 0 {
				list.PopRight()
			}
		}
	}
}

// ----------------------------------------------------------------------------
// Benchmarks: Different Data Types
// ----------------------------------------------------------------------------

func BenchmarkLinkedList_String_Push(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list := NewLinkedList[string]()
		for j := 0; j < 100; j++ {
			list.Push(fmt.Sprintf("string_%d", j))
		}
	}
}

func BenchmarkLinkedList_Struct_Push(b *testing.B) {
	type Data struct {
		ID    int
		Name  string
		Value float64
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list := NewLinkedList[Data]()
		for j := 0; j < 100; j++ {
			list.Push(Data{
				ID:    j,
				Name:  fmt.Sprintf("item_%d", j),
				Value: float64(j) * 1.5,
			})
		}
	}
}

func BenchmarkLinkedList_Pointer_Push(b *testing.B) {
	type Data struct {
		ID    int
		Value int
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list := NewLinkedList[*Data]()
		for j := 0; j < 100; j++ {
			list.Push(&Data{ID: j, Value: j * 2})
		}
	}
}

// ----------------------------------------------------------------------------
// Benchmarks: Memory Allocation Patterns
// ----------------------------------------------------------------------------

func BenchmarkLinkedList_GrowShrinkCycle(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list := NewLinkedList[int]()

		// Grow
		for j := 0; j < 1000; j++ {
			list.Push(j)
		}

		// Shrink
		for j := 0; j < 900; j++ {
			list.PopRight()
		}

		// Grow again
		for j := 0; j < 900; j++ {
			list.Push(j)
		}

		// Final cleanup
		list.DeleteAll()
	}
}

func BenchmarkLinkedList_AlternatingEnds(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list := NewLinkedList[int]()

		for j := 0; j < 500; j++ {
			if j%2 == 0 {
				list.Push(j)
			} else {
				list.PushFront(j)
			}
		}

		for j := 0; j < 500; j++ {
			if j%2 == 0 {
				list.PopRight()
			} else {
				list.PopLeft()
			}
		}
	}
}

// ----------------------------------------------------------------------------
// Benchmarks: First and Last Access
// ----------------------------------------------------------------------------

func BenchmarkLinkedList_First_Repeated(b *testing.B) {
	list := NewLinkedList[int]()
	for i := 0; i < 1000; i++ {
		list.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = list.First()
	}
}

func BenchmarkLinkedList_Last_Repeated(b *testing.B) {
	list := NewLinkedList[int]()
	for i := 0; i < 1000; i++ {
		list.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = list.Last()
	}
}

// ----------------------------------------------------------------------------
// Benchmarks: Size Checks
// ----------------------------------------------------------------------------

func BenchmarkLinkedList_Size_Repeated(b *testing.B) {
	list := NewLinkedList[int]()
	for i := 0; i < 1000; i++ {
		list.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = list.Size()
	}
}

func BenchmarkLinkedList_IsEmpty_Repeated(b *testing.B) {
	list := NewLinkedList[int]()
	for i := 0; i < 1000; i++ {
		list.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = list.IsEmpty()
	}
}
