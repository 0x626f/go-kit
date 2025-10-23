package array

import (
	"github.com/0x626f/go-kit/pkg/number"
	"math/rand"
	"testing"
	"time"
)

func TestInsertionSort(t *testing.T) {
	size := 1_000
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	array := New[int]()

	for range size {
		rnd := r.Intn(10000)
		array.Push(rnd)
	}

	array.InsertionSort(number.NumericComparator[int])

	for index := 1; index < size; index++ {
		if array.At(index-1) > array.At(index) {
			t.Fatal("wrong insertion sort")
		}
	}
}

func TestHeapSort(t *testing.T) {
	size := 1_000_000
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	array := New[int]()

	for range size {
		rnd := r.Intn(10000)
		array.Push(rnd)
	}

	array.HeapSort(number.NumericComparator[int])

	for index := 1; index < size; index++ {
		if array.At(index-1) > array.At(index) {
			t.Fatal("wrong insertion sort")
		}
	}
}

func TestBinarySearch(t *testing.T) {
	size := 1_000_000
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	array := New[int]()

	for range size {
		rnd := r.Intn(10000)
		array.Push(rnd)
	}

	array.HeapSort(number.NumericComparator[int])

	target := array.At(57)

	result, found := array.BinarySearch(target, number.NumericComparator[int])

	if !found {
		t.Fatal("not found")
	}

	if result != target {
		t.Fatal("wrong search")
	}
}
func BenchmarkInsertionSort(b *testing.B) {
	size := 1_000
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	array := New[int]()

	for range size {
		rnd := r.Intn(10000)
		array.Push(rnd)
	}

	b.StartTimer()
	array.InsertionSort(number.NumericComparator[int])
	b.StopTimer()
	b.ReportAllocs()
}

func BenchmarkHeapSort(b *testing.B) {
	size := 1_000_000
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	array := New[int]()

	for range size {
		rnd := r.Intn(10000)
		array.Push(rnd)
	}

	b.StartTimer()
	array.HeapSort(number.NumericComparator[int])
	b.StopTimer()
	b.ReportAllocs()
}

func BenchmarkBinarySearch(b *testing.B) {
	size := 1_000_000
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	array := New[int]()

	for range size {
		rnd := r.Intn(10000)
		array.Push(rnd)
	}

	array.HeapSort(number.NumericComparator[int])

	target := array.At(57)

	b.StartTimer()
	result, found := array.BinarySearch(target, number.NumericComparator[int])
	b.StopTimer()
	b.ReportAllocs()

	if !found {
		b.Fatal("not found")
	}

	if result != target {
		b.Fatal("wrong search")
	}
}
