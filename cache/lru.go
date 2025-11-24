package cache

import (
	"github.com/0x626f/go-kit/linkedlist"
	"github.com/0x626f/go-kit/shared"
	"github.com/0x626f/go-kit/utils"
)

type LRUCache[K comparable, D any] struct {
	capacity int
	recent   *linkedlist.LinkedList[*shared.Pair[K, D], int]
	data     map[K]*linkedlist.LinkedNode[*shared.Pair[K, D]]
}

func NewLRUCache[K comparable, D any](capacity int) *LRUCache[K, D] {
	return &LRUCache[K, D]{
		capacity: capacity,
		recent:   linkedlist.NewLinkedList[*shared.Pair[K, D]](),
		data:     make(map[K]*linkedlist.LinkedNode[*shared.Pair[K, D]]),
	}
}

func (cache *LRUCache[K, D]) Set(key K, item D) {
	if _, exists := cache.data[key]; exists {
		return
	}

	node := cache.recent.InsertFront(&shared.Pair[K, D]{First: key, Second: item})
	cache.data[key] = node

	if cache.capacity != 0 && cache.recent.Size() > cache.capacity {
		retired := cache.recent.PopRight()
		delete(cache.data, retired.First)
	}
}

func (cache *LRUCache[K, D]) Get(key K) (D, bool) {
	if node, exists := cache.data[key]; exists {
		cache.recent.MoveToFront(node)
		return node.Data.Second, true
	}

	return utils.Zero[D](), false
}

func (cache *LRUCache[K, D]) Delete(key K) bool {
	if node, exists := cache.data[key]; exists {
		cache.recent.Remove(node)
		delete(cache.data, key)
		return true
	}
	return false
}

func (cache *LRUCache[K, D]) Flush(count int) {
	if count < 0 {
		return
	}

	if cache.recent.Size() > count {
		cache.recent.ForEach(func(index int, data *shared.Pair[K, D]) bool {
			if (index + 1) > count {
				delete(cache.data, data.First)
			}
			return true
		})
		cache.recent.Shrink(count)
	}
}
