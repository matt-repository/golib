package common

// MapGetElement ...
type MapGetElement[K comparable, V, E any] func(k K, v V) E

// SliceGetElement ...
type SliceGetElement[S ~[]E1, E1, E2 any] func(s S, i int) E2

// SliceExist ...
type SliceExist[S ~[]E, E any] func(s S, i int) bool

// GetKeyValue get key value
type GetKeyValue[K comparable, V, E any] func(v E) (key K, value V)

// EqualElement ...
type EqualElement[E any] func(v1, v2 E) bool

// Less ...
type Less[S ~[]E, E any] func(s S, i, j int) bool
