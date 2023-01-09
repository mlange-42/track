package util

// Pair is a (key/value) pair
type Pair[K any, V any] struct {
	Key   K
	Value V
}

// NewPair creates a new Pair
func NewPair[K any, V any](key K, value V) Pair[K, V] {
	return Pair[K, V]{key, value}
}
