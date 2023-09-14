package _map

import "github.com/matt-repository/matt_golib/common"

// ToSlice ...
func ToSlice[M ~map[K]V, K comparable, S []E, V any, E any](mapData M, f common.MapGetElement[K, V, E]) S {
	ret := make(S, 0)
	for k, v := range mapData {
		s := f(k, v)
		ret = append(ret, s)
	}
	return ret
}

// Equal ..
func Equal[M ~map[K]V, K comparable, V any](m1, m2 M, equal common.EqualElement[V]) bool {
	if len(m1) != len(m2) {
		return false
	}
	for k, v1 := range m1 {
		if v2, ok := m2[k]; !ok || !equal(v1, v2) {
			return false
		}
	}
	return true
}
