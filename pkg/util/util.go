package util

import (
	"sort"

	"golang.org/x/exp/constraints"
)

// Set[T] is an unordered set of T.
type Set[T comparable] map[T]struct{}

func NewSet[T comparable](values ...T) Set[T] {
	ret := make(Set[T], len(values))
	for _, value := range values {
		ret.Insert(value)
	}
	return ret
}

func (o Set[T]) Insert(v T) {
	o[v] = struct{}{}
}

func MapKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

func SortedMapKeys[K constraints.Ordered, V any](m map[K]V) []K {
	keys := MapKeys(m)
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	return keys
}
