package set

import (
	"testing"
)

type User struct {
	Id int
}

func (user *User) Key() int {
	return user.Id
}

func createSampleUsers(count, startFrom int) []*User {
	var users []*User
	for index := range count {
		users = append(users, &User{Id: startFrom + index})
	}
	return users
}

func TestWrap(t *testing.T) {
	array := createSampleUsers(10, 0)
	set := Wrap(array...)
	for _, item := range array {
		if item.Id != set.At(item.Id).Id {
			t.Fatal("item mismatch")
		}
	}
}

func TestDuplicates(t *testing.T) {
	set := Wrap(append(createSampleUsers(10, 0), createSampleUsers(10, 0)...)...)

	if set.Size() != 10 {
		t.Fatal("duplicates were stored")
	}
}

func TestGet(t *testing.T) {
	set := Wrap(createSampleUsers(10, 0)...)

	item0 := set.Get(0)

	if item0 == nil {
		t.Fatal("item was not found")
	}
}

func TestFilter(t *testing.T) {
	set := Wrap(createSampleUsers(10, 1)...)
	predicate := func(arg *User) bool { return arg.Id%2 == 0 }
	filtered := set.Filter(predicate)

	if filtered.Size() != 5 {
		t.Fatal("wrong size")
	}

}

func TestFind(t *testing.T) {
	set := Wrap(createSampleUsers(10, 0)...)
	target := 5
	predicate := func(arg *User) bool { return arg.Id == target }

	if result, found := set.Find(predicate); !found || result.Id != target {
		t.Fatal("wrong find")
	}
}

func TestJoin(t *testing.T) {
	set := Wrap(createSampleUsers(10, 0)...)
	toJoin := Wrap(createSampleUsers(10, 10)...)

	set.Join(toJoin)

	if set.Size() != 20 {
		t.Fatal("bad join")
	}

	if !set.Has(set.Get(0)) || !set.Has(toJoin.Get(12)) {
		t.Fatal("bad join")
	}
}

func TestMerge(t *testing.T) {
	set0 := Wrap(createSampleUsers(10, 0)...)
	set1 := Wrap(createSampleUsers(10, 10)...)

	merged := set0.Merge(set1).(*Set[int, *User])

	if merged.Size() != 20 {
		t.Fatal("bad join")
	}

	if !merged.Has(set0.Get(0)) || !merged.Has(set1.Get(12)) {
		t.Fatal("bad join")
	}
}

func TestHas(t *testing.T) {
	array := createSampleUsers(10, 0)
	set := Wrap(array...)

	if !set.Has(array[1]) {
		t.Fatal("wrong has")
	}
}

func TestDelete(t *testing.T) {
	set := Wrap(createSampleUsers(10, 0)...)
	sourceSize := set.Size()
	target := 3
	predicate := func(user *User) bool { return user.Id == target }

	set.Delete(target)

	if set.Size() == sourceSize {
		t.Fatal("not shrunk size")
	}

	if set.Some(predicate) {
		t.Fatal("not deleted")
	}
}

func TestDeleteBy(t *testing.T) {
	set := Wrap(createSampleUsers(10, 0)...)
	sourceSize := set.Size()
	target := 3
	predicate := func(user *User) bool { return user.Id == target }

	set.DeleteBy(predicate)

	if set.Size() == sourceSize {
		t.Fatal("not shrunk size")
	}

	if set.Some(predicate) {
		t.Fatal("not deleted")
	}
}

func TestDeleteAll(t *testing.T) {
	set := Wrap(createSampleUsers(10, 0)...)
	sourceSize := set.Size()

	set.DeleteAll()

	if set.Size() == sourceSize {
		t.Fatal("not shrunk size")
	}
}
