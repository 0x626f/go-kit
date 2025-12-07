package linkedlist

import (
	"testing"
)

// ============================================================================
// COMPREHENSIVE TEST SUITE FOR LINKED LIST
// ============================================================================

// ----------------------------------------------------------------------------
// Edge Cases: Empty List Operations
// ----------------------------------------------------------------------------

func TestLinkedList_Edge_EmptyList_AllOperations(t *testing.T) {
	list := NewLinkedList[int]()

	// Test all read operations on empty list
	if list.Size() != 0 {
		t.Errorf("Empty list size should be 0, got %d", list.Size())
	}
	if !list.IsEmpty() {
		t.Error("IsEmpty should return true for empty list")
	}
	if val := list.First(); val != 0 {
		t.Errorf("First on empty list should return zero value, got %d", val)
	}
	if val := list.Last(); val != 0 {
		t.Errorf("Last on empty list should return zero value, got %d", val)
	}
	if val := list.At(0); val != 0 {
		t.Errorf("At(0) on empty list should return zero value, got %d", val)
	}
	if val := list.Get(5); val != 0 {
		t.Errorf("Get(5) on empty list should return zero value, got %d", val)
	}

	// Test pop operations on empty list
	if val := list.PopLeft(); val != 0 {
		t.Errorf("PopLeft on empty list should return zero value, got %d", val)
	}
	if val := list.PopRight(); val != 0 {
		t.Errorf("PopRight on empty list should return zero value, got %d", val)
	}
	if val := list.Pop(0); val != 0 {
		t.Errorf("Pop(0) on empty list should return zero value, got %d", val)
	}

	// Test delete operations on empty list (should not panic)
	list.Delete(0)
	list.Delete(-1)
	list.DeleteBy(func(v int) bool { return true })
	list.DeleteAll()

	// Test search operations on empty list
	if idx, found := list.IndexOf(func(v int) bool { return v == 5 }); found {
		t.Errorf("IndexOf on empty list should not find anything, got index %d", idx)
	}
	if val, found := list.Find(func(v int) bool { return true }); found {
		t.Errorf("Find on empty list should not find anything, got %d", val)
	}
	if list.Some(func(v int) bool { return true }) {
		t.Error("Some on empty list should return false")
	}
}

// ----------------------------------------------------------------------------
// Edge Cases: Single Element Operations
// ----------------------------------------------------------------------------

func TestLinkedList_Edge_SingleElement_Comprehensive(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func() *LinkedList[int]
	}{
		{"Push", func() *LinkedList[int] {
			l := NewLinkedList[int]()
			l.Push(42)
			return l
		}},
		{"PushFront", func() *LinkedList[int] {
			l := NewLinkedList[int]()
			l.PushFront(42)
			return l
		}},
		{"Insert", func() *LinkedList[int] {
			l := NewLinkedList[int]()
			l.Insert(42)
			return l
		}},
		{"InsertFront", func() *LinkedList[int] {
			l := NewLinkedList[int]()
			l.InsertFront(42)
			return l
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := tt.setupFunc()

			// Verify basic properties
			if list.Size() != 1 {
				t.Errorf("Size should be 1, got %d", list.Size())
			}
			if list.IsEmpty() {
				t.Error("IsEmpty should return false")
			}

			// Test access methods
			if val := list.First(); val != 42 {
				t.Errorf("First should return 42, got %d", val)
			}
			if val := list.Last(); val != 42 {
				t.Errorf("Last should return 42, got %d", val)
			}
			if val := list.At(0); val != 42 {
				t.Errorf("At(0) should return 42, got %d", val)
			}
			if val := list.At(-1); val != 42 {
				t.Errorf("At(-1) should return 42, got %d", val)
			}

			// Test out of bounds
			if val := list.At(1); val != 0 {
				t.Errorf("At(1) out of bounds should return 0, got %d", val)
			}
			if val := list.At(-2); val != 0 {
				t.Errorf("At(-2) out of bounds should return 0, got %d", val)
			}

			// Test PopLeft (should empty the list)
			list2 := tt.setupFunc()
			if val := list2.PopLeft(); val != 42 {
				t.Errorf("PopLeft should return 42, got %d", val)
			}
			if list2.Size() != 0 {
				t.Errorf("Size after PopLeft should be 0, got %d", list2.Size())
			}

			// Test PopRight (should empty the list)
			list3 := tt.setupFunc()
			if val := list3.PopRight(); val != 42 {
				t.Errorf("PopRight should return 42, got %d", val)
			}
			if list3.Size() != 0 {
				t.Errorf("Size after PopRight should be 0, got %d", list3.Size())
			}

			// Test Pop with index (should empty the list)
			list4 := tt.setupFunc()
			if val := list4.Pop(0); val != 42 {
				t.Errorf("Pop(0) should return 42, got %d", val)
			}
			if list4.Size() != 0 {
				t.Errorf("Size after Pop(0) should be 0, got %d", list4.Size())
			}
		})
	}
}

// ----------------------------------------------------------------------------
// Edge Cases: Two Element Operations
// ----------------------------------------------------------------------------

func TestLinkedList_Edge_TwoElements_AllPermutations(t *testing.T) {
	// Test all possible orderings and operations
	t.Run("Push_Push", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.Push(1)
		list.Push(2)
		verifyTwoElements(t, list, 1, 2)
	})

	t.Run("PushFront_PushFront", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.PushFront(2)
		list.PushFront(1)
		verifyTwoElements(t, list, 1, 2)
	})

	t.Run("Push_PushFront", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.Push(2)
		list.PushFront(1)
		verifyTwoElements(t, list, 1, 2)
	})

	t.Run("PushFront_Push", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.PushFront(1)
		list.Push(2)
		verifyTwoElements(t, list, 1, 2)
	})
}

func verifyTwoElements(t *testing.T, list *LinkedList[int], first, second int) {
	if list.Size() != 2 {
		t.Errorf("Size should be 2, got %d", list.Size())
	}
	if list.First() != first {
		t.Errorf("First should be %d, got %d", first, list.First())
	}
	if list.Last() != second {
		t.Errorf("Last should be %d, got %d", second, list.Last())
	}
	if list.At(0) != first {
		t.Errorf("At(0) should be %d, got %d", first, list.At(0))
	}
	if list.At(1) != second {
		t.Errorf("At(1) should be %d, got %d", second, list.At(1))
	}
	if list.At(-1) != second {
		t.Errorf("At(-1) should be %d, got %d", second, list.At(-1))
	}
	if list.At(-2) != first {
		t.Errorf("At(-2) should be %d, got %d", first, list.At(-2))
	}
}

func TestLinkedList_Edge_TwoElements_PopOperations(t *testing.T) {
	// PopLeft on 2 elements
	t.Run("PopLeft", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.PushAll(1, 2)
		val := list.PopLeft()
		if val != 1 {
			t.Errorf("PopLeft should return 1, got %d", val)
		}
		if list.Size() != 1 {
			t.Errorf("Size should be 1, got %d", list.Size())
		}
		if list.First() != 2 || list.Last() != 2 {
			t.Error("Remaining element should be 2")
		}
	})

	// PopRight on 2 elements
	t.Run("PopRight", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.PushAll(1, 2)
		val := list.PopRight()
		if val != 2 {
			t.Errorf("PopRight should return 2, got %d", val)
		}
		if list.Size() != 1 {
			t.Errorf("Size should be 1, got %d", list.Size())
		}
		if list.First() != 1 || list.Last() != 1 {
			t.Error("Remaining element should be 1")
		}
	})

	// Pop first element
	t.Run("Pop(0)", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.PushAll(1, 2)
		val := list.Pop(0)
		if val != 1 {
			t.Errorf("Pop(0) should return 1, got %d", val)
		}
		if list.Size() != 1 {
			t.Errorf("Size should be 1, got %d", list.Size())
		}
		if list.First() != 2 {
			t.Error("Remaining element should be 2")
		}
	})

	// Pop second element
	t.Run("Pop(1)", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.PushAll(1, 2)
		val := list.Pop(1)
		if val != 2 {
			t.Errorf("Pop(1) should return 2, got %d", val)
		}
		if list.Size() != 1 {
			t.Errorf("Size should be 1, got %d", list.Size())
		}
		if list.First() != 1 {
			t.Error("Remaining element should be 1")
		}
	})
}

// ----------------------------------------------------------------------------
// Edge Cases: Negative Indices
// ----------------------------------------------------------------------------

func TestLinkedList_Edge_NegativeIndices_Comprehensive(t *testing.T) {
	list := NewLinkedList[int]()
	list.PushAll(10, 20, 30, 40, 50)

	tests := []struct {
		index    int
		expected int
	}{
		{-1, 50},
		{-2, 40},
		{-3, 30},
		{-4, 20},
		{-5, 10},
		{-6, 0}, // Out of bounds
		{-100, 0},
	}

	for _, tt := range tests {
		t.Run("At", func(t *testing.T) {
			if val := list.At(tt.index); val != tt.expected {
				t.Errorf("At(%d) expected %d, got %d", tt.index, tt.expected, val)
			}
		})
	}

	// Test Pop with negative indices
	list2 := NewLinkedList[int]()
	list2.PushAll(10, 20, 30, 40, 50)
	if val := list2.Pop(-1); val != 50 {
		t.Errorf("Pop(-1) expected 50, got %d", val)
	}
	if list2.Size() != 4 {
		t.Errorf("Size after Pop(-1) should be 4, got %d", list2.Size())
	}

	if val := list2.Pop(-2); val != 30 {
		t.Errorf("Pop(-2) expected 30, got %d", val)
	}
	if list2.Size() != 3 {
		t.Errorf("Size after Pop(-2) should be 3, got %d", list2.Size())
	}
}

// ----------------------------------------------------------------------------
// Edge Cases: Swap Operations
// ----------------------------------------------------------------------------

func TestLinkedList_Edge_Swap_AllCombinations(t *testing.T) {
	// Swap with same index (no-op)
	t.Run("SameIndex", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.PushAll(1, 2, 3, 4, 5)
		list.Swap(2, 2)
		expected := []int{1, 2, 3, 4, 5}
		verifySequence(t, list, expected)
	})

	// Swap adjacent elements
	t.Run("Adjacent", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.PushAll(1, 2, 3, 4, 5)
		list.Swap(1, 2)
		expected := []int{1, 3, 2, 4, 5}
		verifySequence(t, list, expected)
	})

	// Swap first and second
	t.Run("FirstAndSecond", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.PushAll(1, 2, 3, 4, 5)
		list.Swap(0, 1)
		expected := []int{2, 1, 3, 4, 5}
		verifySequence(t, list, expected)
	})

	// Swap last two
	t.Run("LastTwo", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.PushAll(1, 2, 3, 4, 5)
		list.Swap(3, 4)
		expected := []int{1, 2, 3, 5, 4}
		verifySequence(t, list, expected)
	})

	// Swap first and last
	t.Run("FirstAndLast", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.PushAll(1, 2, 3, 4, 5)
		list.Swap(0, 4)
		expected := []int{5, 2, 3, 4, 1}
		verifySequence(t, list, expected)
	})

	// Swap with negative indices
	t.Run("NegativeIndices", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.PushAll(1, 2, 3, 4, 5)
		list.Swap(-1, -2)
		expected := []int{1, 2, 3, 5, 4}
		verifySequence(t, list, expected)
	})

	// Swap distant elements
	t.Run("Distant", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.PushAll(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
		list.Swap(1, 8)
		expected := []int{1, 9, 3, 4, 5, 6, 7, 8, 2, 10}
		verifySequence(t, list, expected)
	})

	// Swap out of bounds (should be no-op)
	t.Run("OutOfBounds", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.PushAll(1, 2, 3)
		list.Swap(0, 10)
		expected := []int{1, 2, 3}
		verifySequence(t, list, expected)
	})

	// Swap on two-element list
	t.Run("TwoElements", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.PushAll(1, 2)
		list.Swap(0, 1)
		expected := []int{2, 1}
		verifySequence(t, list, expected)
	})
}

// ----------------------------------------------------------------------------
// Edge Cases: Move Operations
// ----------------------------------------------------------------------------

func TestLinkedList_Edge_Move_Comprehensive(t *testing.T) {
	// Move to same position (no-op)
	t.Run("SamePosition", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.PushAll(1, 2, 3, 4, 5)
		list.Move(2, 2)
		expected := []int{1, 2, 3, 4, 5}
		verifySequence(t, list, expected)
	})

	// Move forward
	t.Run("MoveForward", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.PushAll(1, 2, 3, 4, 5)
		list.Move(1, 3) // Move 2 to position of 4
		expected := []int{1, 3, 4, 2, 5}
		verifySequence(t, list, expected)
	})

	// Move backward
	t.Run("MoveBackward", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.PushAll(1, 2, 3, 4, 5)
		list.Move(3, 1) // Move 4 to position of 2
		expected := []int{1, 4, 2, 3, 5}
		verifySequence(t, list, expected)
	})

	// Move first to last
	t.Run("FirstToLast", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.PushAll(1, 2, 3, 4, 5)
		list.Move(0, 4)
		expected := []int{2, 3, 4, 5, 1}
		verifySequence(t, list, expected)
	})

	// Move last to first
	t.Run("LastToFirst", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.PushAll(1, 2, 3, 4, 5)
		list.Move(4, 0)
		expected := []int{5, 1, 2, 3, 4}
		verifySequence(t, list, expected)
	})

	// Move adjacent elements
	t.Run("Adjacent", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.PushAll(1, 2, 3, 4, 5)
		list.Move(1, 2)
		if list.Size() != 5 {
			t.Errorf("Size should remain 5, got %d", list.Size())
		}
	})
}

// ----------------------------------------------------------------------------
// Edge Cases: Delete Operations
// ----------------------------------------------------------------------------

func TestLinkedList_Edge_Delete_Comprehensive(t *testing.T) {
	// Delete first element
	t.Run("DeleteFirst", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.PushAll(1, 2, 3, 4, 5)
		list.Delete(0)
		expected := []int{2, 3, 4, 5}
		verifySequence(t, list, expected)
	})

	// Delete last element
	t.Run("DeleteLast", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.PushAll(1, 2, 3, 4, 5)
		list.Delete(4)
		expected := []int{1, 2, 3, 4}
		verifySequence(t, list, expected)
	})

	// Delete middle element
	t.Run("DeleteMiddle", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.PushAll(1, 2, 3, 4, 5)
		list.Delete(2)
		expected := []int{1, 2, 4, 5}
		verifySequence(t, list, expected)
	})

	// Delete with negative index
	t.Run("DeleteNegativeIndex", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.PushAll(1, 2, 3, 4, 5)
		list.Delete(-1)
		expected := []int{1, 2, 3, 4}
		verifySequence(t, list, expected)
	})

	// Delete out of bounds (should be no-op)
	t.Run("DeleteOutOfBounds", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.PushAll(1, 2, 3)
		list.Delete(10)
		expected := []int{1, 2, 3}
		verifySequence(t, list, expected)
	})

	// Delete all matching elements
	t.Run("DeleteByMultiple", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.PushAll(1, 2, 3, 2, 4, 2, 5)
		list.DeleteBy(func(v int) bool { return v == 2 })
		expected := []int{1, 3, 4, 5}
		verifySequence(t, list, expected)
	})

	// Delete all elements
	t.Run("DeleteAll", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.PushAll(1, 2, 3, 4, 5)
		list.DeleteAll()
		if list.Size() != 0 {
			t.Errorf("Size should be 0 after DeleteAll, got %d", list.Size())
		}
		if !list.IsEmpty() {
			t.Error("List should be empty after DeleteAll")
		}
	})
}

// ----------------------------------------------------------------------------
// Edge Cases: Shrink Operations
// ----------------------------------------------------------------------------

func TestLinkedList_Edge_Shrink_Comprehensive(t *testing.T) {
	// Shrink to 0
	t.Run("ShrinkToZero", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.PushAll(1, 2, 3, 4, 5)
		list.Shrink(0)
		if list.Size() != 0 {
			t.Errorf("Size should be 0, got %d", list.Size())
		}
		if !list.IsEmpty() {
			t.Error("List should be empty")
		}
	})

	// Shrink to 1
	t.Run("ShrinkToOne", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.PushAll(1, 2, 3, 4, 5)
		list.Shrink(1)
		if list.Size() != 1 {
			t.Errorf("Size should be 1, got %d", list.Size())
		}
		if list.First() != 1 {
			t.Errorf("First element should be 1, got %d", list.First())
		}
	})

	// Shrink to middle
	t.Run("ShrinkToMiddle", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.PushAll(1, 2, 3, 4, 5)
		list.Shrink(3)
		expected := []int{1, 2, 3}
		verifySequence(t, list, expected)
	})

	// Shrink larger than size (no-op)
	t.Run("ShrinkLargerThanSize", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.PushAll(1, 2, 3)
		list.Shrink(10)
		expected := []int{1, 2, 3}
		verifySequence(t, list, expected)
	})

	// Shrink equal to size (no-op)
	t.Run("ShrinkEqualToSize", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.PushAll(1, 2, 3)
		list.Shrink(3)
		expected := []int{1, 2, 3}
		verifySequence(t, list, expected)
	})
}

// ----------------------------------------------------------------------------
// Edge Cases: Join and Merge Operations
// ----------------------------------------------------------------------------

func TestLinkedList_Edge_JoinMerge_Comprehensive(t *testing.T) {
	// Join empty lists
	t.Run("JoinBothEmpty", func(t *testing.T) {
		list1 := NewLinkedList[int]()
		list2 := NewLinkedList[int]()
		list1.Join(list2)
		if list1.Size() != 0 {
			t.Errorf("Size should be 0, got %d", list1.Size())
		}
	})

	t.Run("JoinFirstEmpty", func(t *testing.T) {
		list1 := NewLinkedList[int]()
		list2 := NewLinkedList[int]()
		list2.PushAll(1, 2, 3)
		list1.Join(list2)
		expected := []int{1, 2, 3}
		verifySequence(t, list1, expected)
	})

	t.Run("JoinSecondEmpty", func(t *testing.T) {
		list1 := NewLinkedList[int]()
		list1.PushAll(1, 2, 3)
		list2 := NewLinkedList[int]()
		list1.Join(list2)
		expected := []int{1, 2, 3}
		verifySequence(t, list1, expected)
	})

	// Merge operations
	t.Run("MergeBothEmpty", func(t *testing.T) {
		list1 := NewLinkedList[int]()
		list2 := NewLinkedList[int]()
		merged := list1.Merge(list2)
		if merged.Size() != 0 {
			t.Errorf("Merged size should be 0, got %d", merged.Size())
		}
	})

	t.Run("MergePreservesOriginal", func(t *testing.T) {
		list1 := NewLinkedList[int]()
		list1.PushAll(1, 2, 3)
		list2 := NewLinkedList[int]()
		list2.PushAll(4, 5, 6)
		merged := list1.Merge(list2)

		// Check merged
		if merged.Size() != 6 {
			t.Errorf("Merged size should be 6, got %d", merged.Size())
		}

		// Check originals are unchanged
		if list1.Size() != 3 {
			t.Errorf("Original list1 size should be 3, got %d", list1.Size())
		}
		if list2.Size() != 3 {
			t.Errorf("Original list2 size should be 3, got %d", list2.Size())
		}
	})
}

// ----------------------------------------------------------------------------
// Edge Cases: Complex Type Operations
// ----------------------------------------------------------------------------

func TestLinkedList_Edge_ComplexTypes(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	// Test with struct
	t.Run("Struct", func(t *testing.T) {
		list := NewLinkedList[Person]()
		list.Push(Person{"Alice", 30})
		list.Push(Person{"Bob", 25})
		list.Push(Person{"Charlie", 35})

		if list.Size() != 3 {
			t.Errorf("Size should be 3, got %d", list.Size())
		}

		found, exists := list.Find(func(p Person) bool {
			return p.Name == "Bob"
		})
		if !exists {
			t.Error("Should find Bob")
		}
		if found.Age != 25 {
			t.Errorf("Bob's age should be 25, got %d", found.Age)
		}
	})

	// Test with pointer
	t.Run("Pointer", func(t *testing.T) {
		list := NewLinkedList[*Person]()
		p1 := &Person{"Alice", 30}
		p2 := &Person{"Bob", 25}
		list.Push(p1)
		list.Push(p2)

		// Modify through pointer
		p1.Age = 31
		if list.First().Age != 31 {
			t.Errorf("Age should be 31, got %d", list.First().Age)
		}
	})

	// Test with slice
	t.Run("Slice", func(t *testing.T) {
		list := NewLinkedList[[]int]()
		list.Push([]int{1, 2, 3})
		list.Push([]int{4, 5, 6})

		if len(list.First()) != 3 {
			t.Errorf("First slice length should be 3, got %d", len(list.First()))
		}
	})

	// Test with map
	t.Run("Map", func(t *testing.T) {
		list := NewLinkedList[map[string]int]()
		list.Push(map[string]int{"a": 1, "b": 2})
		list.Push(map[string]int{"c": 3, "d": 4})

		if len(list.First()) != 2 {
			t.Errorf("First map length should be 2, got %d", len(list.First()))
		}
	})
}

// ----------------------------------------------------------------------------
// Edge Cases: Stress Tests
// ----------------------------------------------------------------------------

func TestLinkedList_Edge_Stress_RandomOperations(t *testing.T) {
	list := NewLinkedList[int]()

	// Add 1000 elements
	for i := 0; i < 1000; i++ {
		if i%2 == 0 {
			list.Push(i)
		} else {
			list.PushFront(i)
		}
	}

	if list.Size() != 1000 {
		t.Errorf("Size should be 1000, got %d", list.Size())
	}

	// Delete every 3rd element
	count := 0
	for i := 0; i < list.Size(); {
		if count%3 == 0 {
			list.Delete(i)
		} else {
			i++
		}
		count++
	}

	// Verify list is still functional
	if list.IsEmpty() {
		t.Error("List should not be empty")
	}

	// Test access
	_ = list.First()
	_ = list.Last()
	_ = list.At(list.Size() / 2)
}

func TestLinkedList_Edge_Stress_AlternatingPopPush(t *testing.T) {
	list := NewLinkedList[int]()

	// Perform 10000 alternating push/pop operations
	for i := 0; i < 10000; i++ {
		list.Push(i)
		if i%2 == 0 && list.Size() > 1 {
			list.PopLeft()
		}
		if i%3 == 0 && list.Size() > 1 {
			list.PopRight()
		}
	}

	// Verify list is still functional
	if list.IsEmpty() {
		t.Error("List should not be empty after stress test")
	}

	// Pop all remaining elements
	initialSize := list.Size()
	for i := 0; i < initialSize; i++ {
		if i%2 == 0 {
			list.PopLeft()
		} else {
			list.PopRight()
		}
	}

	if list.Size() != 0 {
		t.Errorf("Size should be 0 after popping all elements, got %d", list.Size())
	}
}

// ----------------------------------------------------------------------------
// Edge Cases: Filter and ForEach
// ----------------------------------------------------------------------------

func TestLinkedList_Edge_FilterForEach_Comprehensive(t *testing.T) {
	// Filter all elements
	t.Run("FilterAll", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.PushAll(1, 2, 3, 4, 5)
		filtered := list.Filter(func(v int) bool { return true })
		if filtered.Size() != 5 {
			t.Errorf("Filtered size should be 5, got %d", filtered.Size())
		}
	})

	// Filter none
	t.Run("FilterNone", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.PushAll(1, 2, 3, 4, 5)
		filtered := list.Filter(func(v int) bool { return false })
		if filtered.Size() != 0 {
			t.Errorf("Filtered size should be 0, got %d", filtered.Size())
		}
	})

	// ForEach with early exit
	t.Run("ForEachEarlyExit", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.PushAll(1, 2, 3, 4, 5)
		count := 0
		list.ForEach(func(idx int, val int) bool {
			count++
			return val < 3
		})
		if count != 3 {
			t.Errorf("Should have visited 3 elements, visited %d", count)
		}
	})

	// ForEach on empty list
	t.Run("ForEachEmpty", func(t *testing.T) {
		list := NewLinkedList[int]()
		count := 0
		list.ForEach(func(idx int, val int) bool {
			count++
			return true
		})
		if count != 0 {
			t.Errorf("Should not have visited any elements, visited %d", count)
		}
	})
}

// ----------------------------------------------------------------------------
// Edge Cases: MoveToFront Operations
// ----------------------------------------------------------------------------

func TestLinkedList_MoveToFront_MiddleNode(t *testing.T) {
	list := NewLinkedList[int]()
	list.PushAll(1, 2, 3, 4, 5)

	// Get the node at index 2 (value 3)
	node := list.findNodeByIndex(2)
	if node == nil {
		t.Fatal("Node should not be nil")
	}

	list.MoveToFront(node)

	// Expected: 3, 1, 2, 4, 5
	expected := []int{3, 1, 2, 4, 5}
	verifySequence(t, list, expected)

	if list.First() != 3 {
		t.Errorf("First element should be 3, got %d", list.First())
	}
}

func TestLinkedList_MoveToFront_LastNode(t *testing.T) {
	list := NewLinkedList[int]()
	list.PushAll(1, 2, 3, 4, 5)

	// Get the last node (value 5)
	node := list.findNodeByIndex(4)
	if node == nil {
		t.Fatal("Node should not be nil")
	}

	list.MoveToFront(node)

	// Expected: 5, 1, 2, 3, 4
	expected := []int{5, 1, 2, 3, 4}
	verifySequence(t, list, expected)

	if list.First() != 5 {
		t.Errorf("First element should be 5, got %d", list.First())
	}
	if list.Last() != 4 {
		t.Errorf("Last element should be 4, got %d", list.Last())
	}
}

func TestLinkedList_MoveToFront_HeadNode(t *testing.T) {
	list := NewLinkedList[int]()
	list.PushAll(1, 2, 3, 4, 5)

	// Get the head node (value 1)
	node := list.findNodeByIndex(0)
	if node == nil {
		t.Fatal("Node should not be nil")
	}

	list.MoveToFront(node)

	// Expected: no change (1, 2, 3, 4, 5)
	expected := []int{1, 2, 3, 4, 5}
	verifySequence(t, list, expected)

	if list.First() != 1 {
		t.Errorf("First element should still be 1, got %d", list.First())
	}
}

func TestLinkedList_MoveToFront_SecondNode(t *testing.T) {
	list := NewLinkedList[int]()
	list.PushAll(1, 2, 3, 4, 5)

	// Get the second node (value 2)
	node := list.findNodeByIndex(1)
	if node == nil {
		t.Fatal("Node should not be nil")
	}

	list.MoveToFront(node)

	// Expected: 2, 1, 3, 4, 5
	expected := []int{2, 1, 3, 4, 5}
	verifySequence(t, list, expected)
}

func TestLinkedList_MoveToFront_TwoElementList(t *testing.T) {
	// Move second to front
	t.Run("MoveSecondToFront", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.PushAll(1, 2)

		node := list.findNodeByIndex(1)
		list.MoveToFront(node)

		expected := []int{2, 1}
		verifySequence(t, list, expected)
	})

	// Move first to front (no-op)
	t.Run("MoveFirstToFront", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.PushAll(1, 2)

		node := list.findNodeByIndex(0)
		list.MoveToFront(node)

		expected := []int{1, 2}
		verifySequence(t, list, expected)
	})
}

func TestLinkedList_MoveToFront_SingleElementList(t *testing.T) {
	list := NewLinkedList[int]()
	list.Push(42)

	node := list.findNodeByIndex(0)
	if node == nil {
		t.Fatal("Node should not be nil")
	}

	list.MoveToFront(node)

	// Expected: no change (42)
	expected := []int{42}
	verifySequence(t, list, expected)

	if list.Size() != 1 {
		t.Errorf("Size should still be 1, got %d", list.Size())
	}
}

func TestLinkedList_MoveToFront_NilNode(t *testing.T) {
	list := NewLinkedList[int]()
	list.PushAll(1, 2, 3, 4, 5)

	// Move nil node (should be no-op)
	list.MoveToFront(nil)

	// Expected: no change
	expected := []int{1, 2, 3, 4, 5}
	verifySequence(t, list, expected)
}

func TestLinkedList_MoveToFront_InsertedNode(t *testing.T) {
	list := NewLinkedList[int]()
	list.PushAll(1, 2, 4, 5)

	// Insert a node at the end and get its reference
	node := list.Insert(3)

	// Now move it to front
	list.MoveToFront(node)

	// Expected: 3, 1, 2, 4, 5
	expected := []int{3, 1, 2, 4, 5}
	verifySequence(t, list, expected)

	if list.First() != 3 {
		t.Errorf("First element should be 3, got %d", list.First())
	}
}

func TestLinkedList_MoveToFront_InsertFrontNode(t *testing.T) {
	list := NewLinkedList[int]()
	list.PushAll(2, 3, 4, 5)

	// Insert a node at the front and get its reference
	node := list.InsertFront(1)

	// Node is already at front, move it again (should be no-op)
	list.MoveToFront(node)

	// Expected: no change (1, 2, 3, 4, 5)
	expected := []int{1, 2, 3, 4, 5}
	verifySequence(t, list, expected)
}

func TestLinkedList_MoveToFront_MultipleMoves(t *testing.T) {
	list := NewLinkedList[int]()
	list.PushAll(1, 2, 3, 4, 5)

	// Move node 4 to front
	node4 := list.findNodeByIndex(3)
	list.MoveToFront(node4)
	// Expected: 4, 1, 2, 3, 5
	expected1 := []int{4, 1, 2, 3, 5}
	verifySequence(t, list, expected1)

	// Move node 5 to front
	node5 := list.findNodeByIndex(4)
	list.MoveToFront(node5)
	// Expected: 5, 4, 1, 2, 3
	expected2 := []int{5, 4, 1, 2, 3}
	verifySequence(t, list, expected2)

	// Move node 2 to front
	node2 := list.findNodeByIndex(3)
	list.MoveToFront(node2)
	// Expected: 2, 5, 4, 1, 3
	expected3 := []int{2, 5, 4, 1, 3}
	verifySequence(t, list, expected3)
}

func TestLinkedList_MoveToFront_AfterPop(t *testing.T) {
	list := NewLinkedList[int]()
	list.PushAll(1, 2, 3, 4, 5)

	// Get reference to node 3
	node3 := list.findNodeByIndex(2)

	// Pop first element
	list.PopLeft()
	// List is now: 2, 3, 4, 5

	// Move node 3 to front
	list.MoveToFront(node3)
	// Expected: 3, 2, 4, 5
	expected := []int{3, 2, 4, 5}
	verifySequence(t, list, expected)
}

func TestLinkedList_MoveToFront_LargeList(t *testing.T) {
	list := NewLinkedList[int]()

	// Create a list with 100 elements
	for i := 1; i <= 100; i++ {
		list.Push(i)
	}

	// Move element 50 to front
	node50 := list.findNodeByIndex(49)
	list.MoveToFront(node50)

	if list.First() != 50 {
		t.Errorf("First element should be 50, got %d", list.First())
	}
	if list.Size() != 100 {
		t.Errorf("Size should still be 100, got %d", list.Size())
	}

	// Move element 100 to front
	node100 := list.findNodeByIndex(99)
	list.MoveToFront(node100)

	if list.First() != 100 {
		t.Errorf("First element should be 100, got %d", list.First())
	}
	if list.Size() != 100 {
		t.Errorf("Size should still be 100, got %d", list.Size())
	}
}

func TestLinkedList_MoveToFront_PreservesDataIntegrity(t *testing.T) {
	type Item struct {
		ID   int
		Name string
	}

	list := NewLinkedList[Item]()
	list.Push(Item{1, "First"})
	list.Push(Item{2, "Second"})
	list.Push(Item{3, "Third"})
	list.Push(Item{4, "Fourth"})

	// Get the third node
	node := list.findNodeByIndex(2)

	// Verify Data before move
	if node.Data.ID != 3 || node.Data.Name != "Third" {
		t.Error("Node Data incorrect before move")
	}

	list.MoveToFront(node)

	// Verify Data after move
	if list.First().ID != 3 || list.First().Name != "Third" {
		t.Error("First element Data incorrect after move")
	}

	// Verify the node still has correct Data
	if node.Data.ID != 3 || node.Data.Name != "Third" {
		t.Error("Node Data incorrect after move")
	}
}

func TestLinkedList_MoveToFront_ConsecutiveMovesSameNode(t *testing.T) {
	list := NewLinkedList[int]()
	list.PushAll(1, 2, 3, 4, 5)

	// Get node 3
	node := list.findNodeByIndex(2)

	// Move to front multiple times
	list.MoveToFront(node)
	list.MoveToFront(node)
	list.MoveToFront(node)

	// Expected: 3, 1, 2, 4, 5 (should be same after first move)
	expected := []int{3, 1, 2, 4, 5}
	verifySequence(t, list, expected)
}

// ----------------------------------------------------------------------------
// Sort Operations
// ----------------------------------------------------------------------------

func TestLinkedList_QuickSort_EmptyList(t *testing.T) {
	list := NewLinkedList[int]()

	// Sort empty list (should not panic)
	list.Sort(func(a, b int) int {
		return a - b
	})

	if list.Size() != 0 {
		t.Errorf("Size should still be 0, got %d", list.Size())
	}
}

func TestLinkedList_QuickSort_SingleElement(t *testing.T) {
	list := NewLinkedList[int]()
	list.Push(42)

	list.Sort(func(a, b int) int {
		return a - b
	})

	expected := []int{42}
	verifySequence(t, list, expected)
}

func TestLinkedList_QuickSort_TwoElements_Ascending(t *testing.T) {
	t.Run("AlreadySorted", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.PushAll(1, 2)

		list.Sort(func(a, b int) int {
			return a - b
		})

		expected := []int{1, 2}
		verifySequence(t, list, expected)
	})

	t.Run("NeedsSorting", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.PushAll(2, 1)

		list.Sort(func(a, b int) int {
			return a - b
		})

		expected := []int{1, 2}
		verifySequence(t, list, expected)
	})
}

func TestLinkedList_QuickSort_AlreadySorted(t *testing.T) {
	list := NewLinkedList[int]()
	list.PushAll(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	list.Sort(func(a, b int) int {
		return a - b
	})

	expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	verifySequence(t, list, expected)
}

func TestLinkedList_QuickSort_ReverseSorted(t *testing.T) {
	list := NewLinkedList[int]()
	list.PushAll(10, 9, 8, 7, 6, 5, 4, 3, 2, 1)

	list.Sort(func(a, b int) int {
		return a - b
	})

	expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	verifySequence(t, list, expected)
}

func TestLinkedList_QuickSort_RandomOrder(t *testing.T) {
	list := NewLinkedList[int]()
	list.PushAll(5, 2, 8, 1, 9, 3, 7, 4, 6, 10)

	list.Sort(func(a, b int) int {
		return a - b
	})

	expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	verifySequence(t, list, expected)
}

func TestLinkedList_QuickSort_WithDuplicates(t *testing.T) {
	list := NewLinkedList[int]()
	list.PushAll(5, 2, 8, 2, 9, 5, 7, 2, 5, 10)

	list.Sort(func(a, b int) int {
		return a - b
	})

	expected := []int{2, 2, 2, 5, 5, 5, 7, 8, 9, 10}
	verifySequence(t, list, expected)
}

func TestLinkedList_QuickSort_AllEqual(t *testing.T) {
	list := NewLinkedList[int]()
	list.PushAll(5, 5, 5, 5, 5)

	list.Sort(func(a, b int) int {
		return a - b
	})

	expected := []int{5, 5, 5, 5, 5}
	verifySequence(t, list, expected)
}

func TestLinkedList_QuickSort_DescendingOrder(t *testing.T) {
	list := NewLinkedList[int]()
	list.PushAll(5, 2, 8, 1, 9, 3, 7, 4, 6, 10)

	// Sort in descending order
	list.Sort(func(a, b int) int {
		return b - a
	})

	expected := []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
	verifySequence(t, list, expected)
}

func TestLinkedList_QuickSort_NegativeNumbers(t *testing.T) {
	list := NewLinkedList[int]()
	list.PushAll(5, -2, 8, -10, 0, 3, -7, 4, -1, 10)

	list.Sort(func(a, b int) int {
		return a - b
	})

	expected := []int{-10, -7, -2, -1, 0, 3, 4, 5, 8, 10}
	verifySequence(t, list, expected)
}

func TestLinkedList_QuickSort_LargeList(t *testing.T) {
	list := NewLinkedList[int]()

	// Add 100 elements in reverse order
	for i := 100; i >= 1; i-- {
		list.Push(i)
	}

	list.Sort(func(a, b int) int {
		return a - b
	})

	// Verify it's sorted
	if list.Size() != 100 {
		t.Errorf("Size should be 100, got %d", list.Size())
	}

	// Check first few elements
	if list.At(0) != 1 {
		t.Errorf("First element should be 1, got %d", list.At(0))
	}
	if list.At(1) != 2 {
		t.Errorf("Second element should be 2, got %d", list.At(1))
	}

	// Check last few elements
	if list.At(98) != 99 {
		t.Errorf("99th element should be 99, got %d", list.At(98))
	}
	if list.At(99) != 100 {
		t.Errorf("100th element should be 100, got %d", list.At(99))
	}

	// Verify entire sequence is sorted
	prev := 0
	list.ForEach(func(idx int, val int) bool {
		if val < prev {
			t.Errorf("List not sorted at index %d: %d < %d", idx, val, prev)
		}
		prev = val
		return true
	})
}

func TestLinkedList_QuickSort_ComplexType_ByAge(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	list := NewLinkedList[Person]()
	list.Push(Person{"Alice", 30})
	list.Push(Person{"Bob", 25})
	list.Push(Person{"Charlie", 35})
	list.Push(Person{"David", 20})
	list.Push(Person{"Eve", 28})

	// Sort by age
	list.Sort(func(a, b Person) int {
		return a.Age - b.Age
	})

	if list.Size() != 5 {
		t.Errorf("Size should be 5, got %d", list.Size())
	}

	if list.At(0).Name != "David" || list.At(0).Age != 20 {
		t.Errorf("First should be David(20), got %s(%d)", list.At(0).Name, list.At(0).Age)
	}
	if list.At(1).Name != "Bob" || list.At(1).Age != 25 {
		t.Errorf("Second should be Bob(25), got %s(%d)", list.At(1).Name, list.At(1).Age)
	}
	if list.At(2).Name != "Eve" || list.At(2).Age != 28 {
		t.Errorf("Third should be Eve(28), got %s(%d)", list.At(2).Name, list.At(2).Age)
	}
	if list.At(3).Name != "Alice" || list.At(3).Age != 30 {
		t.Errorf("Fourth should be Alice(30), got %s(%d)", list.At(3).Name, list.At(3).Age)
	}
	if list.At(4).Name != "Charlie" || list.At(4).Age != 35 {
		t.Errorf("Fifth should be Charlie(35), got %s(%d)", list.At(4).Name, list.At(4).Age)
	}
}

func TestLinkedList_QuickSort_ComplexType_ByName(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	list := NewLinkedList[Person]()
	list.Push(Person{"Charlie", 35})
	list.Push(Person{"Alice", 30})
	list.Push(Person{"Eve", 28})
	list.Push(Person{"Bob", 25})
	list.Push(Person{"David", 20})

	// Sort by name (lexicographically)
	list.Sort(func(a, b Person) int {
		if a.Name < b.Name {
			return -1
		} else if a.Name > b.Name {
			return 1
		}
		return 0
	})

	if list.Size() != 5 {
		t.Errorf("Size should be 5, got %d", list.Size())
	}

	expectedNames := []string{"Alice", "Bob", "Charlie", "David", "Eve"}
	for i, expectedName := range expectedNames {
		if list.At(i).Name != expectedName {
			t.Errorf("At index %d: expected %s, got %s", i, expectedName, list.At(i).Name)
		}
	}
}

func TestLinkedList_QuickSort_Strings(t *testing.T) {
	list := NewLinkedList[string]()
	list.Push("zebra")
	list.Push("apple")
	list.Push("monkey")
	list.Push("banana")
	list.Push("orange")

	list.Sort(func(a, b string) int {
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	})

	expected := []string{"apple", "banana", "monkey", "orange", "zebra"}
	if list.Size() != len(expected) {
		t.Errorf("Size should be %d, got %d", len(expected), list.Size())
	}

	for i, exp := range expected {
		if list.At(i) != exp {
			t.Errorf("At index %d: expected %s, got %s", i, exp, list.At(i))
		}
	}
}

func TestLinkedList_QuickSort_ThreeElements(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{"Sorted", []int{1, 2, 3}, []int{1, 2, 3}},
		{"Reverse", []int{3, 2, 1}, []int{1, 2, 3}},
		{"FirstLast", []int{2, 1, 3}, []int{1, 2, 3}},
		{"FirstMiddle", []int{2, 3, 1}, []int{1, 2, 3}},
		{"MiddleLast", []int{1, 3, 2}, []int{1, 2, 3}},
		{"LastMiddle", []int{3, 1, 2}, []int{1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewLinkedList[int]()
			list.PushAll(tt.input...)

			list.Sort(func(a, b int) int {
				return a - b
			})

			verifySequence(t, list, tt.expected)
		})
	}
}

func TestLinkedList_QuickSort_PreservesListSize(t *testing.T) {
	list := NewLinkedList[int]()
	list.PushAll(9, 3, 7, 1, 5, 2, 8, 4, 6)

	originalSize := list.Size()

	list.Sort(func(a, b int) int {
		return a - b
	})

	if list.Size() != originalSize {
		t.Errorf("Size should remain %d, got %d", originalSize, list.Size())
	}
}

func TestLinkedList_QuickSort_ManyDuplicates(t *testing.T) {
	list := NewLinkedList[int]()
	list.PushAll(3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5, 8, 9, 7, 9)

	list.Sort(func(a, b int) int {
		return a - b
	})

	expected := []int{1, 1, 2, 3, 3, 4, 5, 5, 5, 6, 7, 8, 9, 9, 9}
	verifySequence(t, list, expected)
}

func TestLinkedList_QuickSort_AfterMultipleOperations(t *testing.T) {
	list := NewLinkedList[int]()

	// Add elements
	list.PushAll(5, 2, 8, 1, 9)

	// Do some operations
	list.Push(3)
	list.PushFront(7)
	list.Delete(2) // Delete element at index 2

	// Sort
	list.Sort(func(a, b int) int {
		return a - b
	})

	// Verify sorted and correct size
	if list.Size() != 6 {
		t.Errorf("Size should be 6, got %d", list.Size())
	}

	// Check it's sorted
	prev := list.At(0)
	for i := 1; i < list.Size(); i++ {
		curr := list.At(i)
		if curr < prev {
			t.Errorf("List not properly sorted at index %d: %d < %d", i, curr, prev)
		}
		prev = curr
	}
}

func TestLinkedList_QuickSort_FloatComparison(t *testing.T) {
	list := NewLinkedList[float64]()
	list.Push(3.14)
	list.Push(1.41)
	list.Push(2.71)
	list.Push(1.73)
	list.Push(0.577)

	list.Sort(func(a, b float64) int {
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	})

	// Verify order
	expected := []float64{0.577, 1.41, 1.73, 2.71, 3.14}
	if list.Size() != len(expected) {
		t.Errorf("Size should be %d, got %d", len(expected), list.Size())
	}

	for i, exp := range expected {
		if list.At(i) != exp {
			t.Errorf("At index %d: expected %f, got %f", i, exp, list.At(i))
		}
	}
}

// ----------------------------------------------------------------------------
// Helper Functions
// ----------------------------------------------------------------------------

func verifySequence(t *testing.T, list *LinkedList[int], expected []int) {
	if list.Size() != len(expected) {
		t.Errorf("Size mismatch: expected %d, got %d", len(expected), list.Size())
		return
	}

	for i, exp := range expected {
		if val := list.At(i); val != exp {
			t.Errorf("At index %d: expected %d, got %d", i, exp, val)
		}
	}

	// Verify with ForEach
	index := 0
	list.ForEach(func(idx int, val int) bool {
		if idx != index {
			t.Errorf("ForEach index mismatch: expected %d, got %d", index, idx)
		}
		if val != expected[index] {
			t.Errorf("ForEach value mismatch at %d: expected %d, got %d", index, expected[index], val)
		}
		index++
		return true
	})
}
