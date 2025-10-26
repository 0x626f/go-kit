package set

import (
	"github.com/0x626f/go-kit/abstract"
	"testing"
)

func TestPrimaryWrap(t *testing.T) {
	array := []int{0, 1, 2, 3, 4, 5}
	set := WrapToPrimarySet(array...)
	for _, item := range array {

		if item != set.At(item).Value() {
			t.Fatal("item mismatch")
		}
	}
}

func TestPrimaryDuplicates(t *testing.T) {
	array := []int{0, 1, 2, 2, 2, 3, 2, 4, 5, 5, 3}
	set := WrapToPrimarySet(array...)

	if set.Size() != 6 {
		t.Fatal("duplicates were stored")
	}
}

func TestPrimaryGet(t *testing.T) {
	array := []int{0, 1, 2, 3, 4, 5}
	set := WrapToPrimarySet(array...)

	item0 := set.Get(5)
	item1 := set.Get(10)

	if item0 == nil {
		t.Fatal("item was not found")
	}

	if item1 != nil {
		t.Fatal("non existing item was found ")
	}
}

func TestPrimaryFilter(t *testing.T) {
	set := WrapToPrimarySet(1, 2, 3, 4, 5)
	predicate := func(arg *abstract.KeyableWrapper[int]) bool { return arg.Value()%2 == 0 }
	filtered := set.Filter(predicate)

	if filtered.Size() != 2 {
		t.Fatal("wrong size")
	}

}

func TestPrimaryFind(t *testing.T) {
	set := WrapToPrimarySet(1, 2, 3, 4, 5)
	target := 5
	predicate := func(arg *abstract.KeyableWrapper[int]) bool { return arg.Value() == target }

	result, found := set.Find(predicate)

	if !found || result.Value() != target {
		t.Fatal("wrong find")
	}
}

func TestPrimaryJoin(t *testing.T) {
	set := WrapToPrimarySet(1, 2, 3, 4, 5)
	toJoin := WrapToPrimarySet(6, 7, 8, 9, 10)

	set.Join(toJoin)

	if set.Size() != 10 {
		t.Fatal("bad join")
	}

	if !set.Has(set.Item(1)) || !set.Has(set.Item(10)) {
		t.Fatal("bad join")
	}
}

func TestPrimaryHas(t *testing.T) {
	set := WrapToPrimarySet(1, 2, 3, 4, 5)
	target := 5

	if !set.Has(set.Item(target)) {
		t.Fatal("wrong has")
	}
}

func TestPrimaryDelete(t *testing.T) {
	set := WrapToPrimarySet(1, 2, 3, 4, 5)
	sourceSize := set.Size()
	target := 3

	set.Delete(target)

	if set.Size() == sourceSize {
		t.Fatal("not shrunk size")
	}

	if set.Has(set.Item(target)) {
		t.Fatal("not deleted")
	}
}

func TestPrimaryDeleteBy(t *testing.T) {
	set := WrapToPrimarySet(1, 2, 3, 4, 5)
	sourceSize := set.Size()
	target := 3
	predicate := func(arg *abstract.KeyableWrapper[int]) bool { return arg.Value() == target }

	set.DeleteBy(predicate)

	if set.Size() == sourceSize {
		t.Fatal("not shrunk size")
	}

	if set.Has(set.Item(target)) {
		t.Fatal("not deleted")
	}
}

func TestPrimaryDeleteAll(t *testing.T) {
	set := WrapToPrimarySet(1, 2, 3, 4, 5)
	sourceSize := set.Size()

	set.DeleteAll()

	if set.Size() == sourceSize {
		t.Fatal("not shrunk size")
	}
}
