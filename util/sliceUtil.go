package util

import (
	"sort"
)

// filter condition
type condition[T any] func(s []T, i int) bool

// Contains determine whether the condition is met in the array
func Contains[T any](data []T, f condition[T]) bool {
	for i := 0; i < len(data); i++ {
		if f(data, i) {
			return true
		}
	}
	return false
}

// Filter ...
func Filter[T any](data []T, f condition[T]) []T {
	r := make([]T, 0)
	for i := 0; i < len(data); i++ {
		if f(data, i) {
			r = append(r, (data)[i])
		}
	}
	return r
}

type getKeyValue[K comparable, V, E any] func(v E) (key K, value V)

// SliceToMap ...
func SliceToMap[M ~map[K]V, K comparable, V, E any](arr []E, f getKeyValue[K, V, E]) M {
	r := make(map[K]V)
	if f == nil {
		return nil
	}
	for i := 0; i < len(arr); i++ {
		key, value := f(arr[i])
		r[key] = value
	}
	return r
}

// EqualSlice ...
func EqualSlice[E any](arr1, arr2 []E, less func(data []E, i, j int) bool, equal func(v1, v2 E) bool) bool {
	if len(arr1) != len(arr2) {
		return false
	}
	if equal == nil {
		return false
	}
	if less != nil {
		less1 := func(i, j int) bool {
			return less(arr1, i, j)
		}
		less2 := func(i, j int) bool {
			return less(arr2, i, j)
		}
		sort.Slice(arr1, less1)
		sort.Slice(arr2, less2)
	}
	for i := 0; i < len(arr1); i++ {
		if !equal(arr1[i], arr2[i]) {
			return false
		}
	}
	return true
}
