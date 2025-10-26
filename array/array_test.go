package array

import (
	"github.com/0x626f/go-kit/number"
	"testing"
)

func TestWrap(t *testing.T) {
	array0 := []int{1, 2, 3, 4, 5}
	array1 := Wrap(1, 2, 3, 4, 5)

	for index, item := range array0 {
		if item != array1.At(index) {
			t.Fatal("item mismatch")
		}
	}
}

func TestGet(t *testing.T) {
	array := Wrap(1, 2, 3, 4, 5)

	if array.Get(0) != 1 {
		t.Fatal("first mismatch")
	}

	if array.Get(10) != 1 {
		t.Fatal("size up overflow failed")
	}

	if array.Get(-10) != 5 {
		t.Fatal("size low overflow failed")
	}

}

func TestFirstAndLast(t *testing.T) {
	array := Wrap(1, 2, 3, 4, 5)

	if array.First() != 1 {
		t.Fatal("wrong first")
	}

	if array.Last() != 5 {
		t.Fatal("wrong last")
	}
}

func TestMinAndMax(t *testing.T) {
	array := Wrap(1, 2, 3, 4, 5)

	if v, found := array.Min(number.NumericComparator[int]); !found || v != 1 {
		t.Fatal("wrong min")
	}

	if v, found := array.Max(number.NumericComparator[int]); !found || v != 5 {
		t.Fatal("wrong max")
	}
}

func TestSlice(t *testing.T) {
	array := Wrap(1, 2, 3, 4, 5)
	slice := array.Slice(1, 4)

	if slice.Size() != 3 {
		t.Fatal("wrong slice size")
	}

	for index := 1; index < 4; index++ {
		if array.At(index) != slice.At(index-1) {
			t.Fatal("wrong items")
		}
	}
}

func TestFilter(t *testing.T) {
	array := Wrap(1, 2, 3, 4, 5)
	filtered := array.Filter(func(arg int) bool { return arg%2 == 0 })

	if filtered.Size() != 2 {
		t.Fatal("wrong size")
	}

	if filtered.At(0) != 2 && filtered.At(1) != 4 {
		t.Fatal("wrong filtering")
	}
}

func TestSome(t *testing.T) {
	array := Wrap(1, 2, 3, 4, 5)
	target := 5

	if !array.Some(func(arg int) bool { return arg == target }) {
		t.Fatal("wrong has")
	}
}

func TestFind(t *testing.T) {
	array := Wrap(1, 2, 3, 4, 5)
	target := 5
	predicate := func(arg int) bool { return arg == target }

	if result, found := array.Find(predicate); !found || result != target {
		t.Fatal("wrong find")
	}
}

func TestJoin(t *testing.T) {
	array := Wrap(1, 2, 3, 4, 5)
	toJoin := Wrap(6, 7, 8, 9, 10)

	array.Join(toJoin)

	if array.Size() != 10 {
		t.Fatal("bad join")
	}

	if array.First() != 1 || array.Last() != 10 {
		t.Fatal("bad order")
	}
}

func TestMerge(t *testing.T) {
	array0 := Wrap(1, 2, 3, 4, 5)
	array1 := Wrap(6, 7, 8, 9, 10)

	merged := array0.Merge(array1).(*Array[int, int])

	if merged.Size() != 10 {
		t.Fatal("bad merge")
	}

	if merged.First() != 1 || merged.Last() != 10 {
		t.Fatal("bad merge")
	}
}

func TestDelete(t *testing.T) {
	array := Wrap(1, 2, 3, 4, 5)
	sourceSize := array.Size()
	target := 3
	index := array.IndexOf(func(arg int) bool { return arg == target })

	array.Delete(index)

	if array.Size() == sourceSize {
		t.Fatal("not shrunk size")
	}

	if array.Some(func(arg int) bool { return arg == target }) {
		t.Fatal("not deleted")
	}
}

func TestDeleteBy(t *testing.T) {
	array := Wrap(1, 2, 3, 4, 5)
	sourceSize := array.Size()
	predicate := func(arg int) bool { return arg%2 != 0 }

	array.DeleteBy(predicate)

	if array.Size() == sourceSize {
		t.Fatal("not shrunk size")
	}

	if array.Some(predicate) {
		t.Fatal("not deleted")
	}
}

func TestDeleteKeepOrdering(t *testing.T) {
	array := Wrap(1, 2, 3, 4, 5)
	sourceSize := array.Size()
	target := 3
	index := array.IndexOf(func(arg int) bool { return arg == target })

	array.DeleteKeepOrdering(index, true)

	if array.Size() == sourceSize {
		t.Fatal("not shrunk size")
	}

	if array.Some(func(arg int) bool { return arg == target }) {
		t.Fatal("not deleted")
	}

	if !array.IsSorted(number.NumericComparator[int]) {
		t.Fatal("inconsistent order")
	}
}

func TestDeleteByKeepOrdering(t *testing.T) {
	array := Wrap(1, 2, 3, 4, 5)
	sourceSize := array.Size()
	predicate := func(arg int) bool { return arg%2 != 0 }

	array.DeleteByKeepOrdering(predicate, true)

	if array.Size() == sourceSize {
		t.Fatal("not shrunk size")
	}

	if array.Some(predicate) {
		t.Fatal("not deleted")
	}

	if !array.IsSorted(number.NumericComparator[int]) {
		t.Fatal("inconsistent order")
	}

}
