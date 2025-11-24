package cache

import (
	"github.com/0x626f/go-kit/linkedlist"
	"github.com/0x626f/go-kit/shared"
	"github.com/0x626f/go-kit/utils"
)

type LFUCache[K comparable, D any] struct {
	capacity int

	frequencies *linkedlist.LinkedList[*shared.Pair[uint, PrimaryCache[K, D]], int]
	data        PrimaryCache[uint, *linkedlist.LinkedNode[*shared.Pair[uint, PrimaryCache[K, D]]]]
	spot        PrimaryCache[K, *linkedlist.LinkedNode[*shared.Pair[uint, PrimaryCache[K, D]]]]
}

func NewLFUCache[K comparable, D any](capacity int) *LFUCache[K, D] {
	return &LFUCache[K, D]{
		capacity:    capacity,
		frequencies: linkedlist.NewLinkedList[*shared.Pair[uint, PrimaryCache[K, D]]](),
		data:        make(PrimaryCache[uint, *linkedlist.LinkedNode[*shared.Pair[uint, PrimaryCache[K, D]]]]),
		spot:        make(PrimaryCache[K, *linkedlist.LinkedNode[*shared.Pair[uint, PrimaryCache[K, D]]]]),
	}
}

func (cache *LFUCache[K, D]) record(freq uint) *linkedlist.LinkedNode[*shared.Pair[uint, PrimaryCache[K, D]]] {
	if cache.data[freq] == nil {
		entry := &shared.Pair[uint, PrimaryCache[K, D]]{First: freq, Second: make(PrimaryCache[K, D])}
		node := cache.frequencies.Insert(entry)
		cache.data[freq] = node
	}
	return cache.data[freq]
}

func (cache *LFUCache[K, D]) Set(key K, item D) {
	if _, exists := cache.spot[key]; !exists {
		node := cache.record(1)
		node.Data.Second[key] = item

		cache.spot[key] = node
	}
}

func (cache *LFUCache[K, D]) Get(key K) (D, bool) {
	node, exists := cache.spot[key]

	if !exists {
		return utils.Zero[D](), false
	}

	nextNode := cache.record(node.Data.First + 1)

	nextNode.Data.Second[key] = node.Data.Second[key]

	delete(node.Data.Second, key)
	cache.spot[key] = nextNode

	return nextNode.Data.Second[key], true
}

func (cache *LFUCache[K, D]) Delete(key K) bool {
	node, exists := cache.spot[key]

	if !exists {
		return true
	}

	delete(node.Data.Second, key)
	delete(cache.spot, key)

	return true
}

func (cache *LFUCache[K, D]) Flush(count int) {
	if count < 0 {
		count = cache.capacity
	}

	if cache.frequencies.Size() > cache.capacity {
		cache.frequencies.Sort(func(arg0, arg1 *shared.Pair[uint, PrimaryCache[K, D]]) int {
			return int(arg1.First) - int(arg0.First)
		})
		cache.frequencies.ForEach(func(index int, data *shared.Pair[uint, PrimaryCache[K, D]]) bool {
			if (index + 1) > cache.capacity {
				for key := range data.Second {
					delete(cache.spot, key)
				}
				delete(cache.data, data.First)
			}
			return true
		})
		cache.frequencies.Shrink(cache.capacity)
	}
}
