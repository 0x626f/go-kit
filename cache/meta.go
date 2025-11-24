package cache

type Cache[D any, K comparable] interface {
	Set(key K, data D)
	Get(key K) (D, bool)
	Delete(key K) bool
}

type PrimaryCache[K comparable, D any] map[K]D
