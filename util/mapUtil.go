package util

type getElement[T any, S, M any] func(k T, v S) M

// MapToSlice ...
func MapToSlice[M ~map[K]V, K, V comparable, S any](mapData M, f getElement[K, V, S]) []S {
	r := make([]S, 0)
	for k, v := range mapData {
		s := f(k, v)
		r = append(r, s)
	}
	return r
}

// EqualMap ..
func EqualMap[M ~map[K]V, K, V comparable](m1, m2 M) bool {
	if len(m1) != len(m2) {
		return false
	}
	for k, v1 := range m1 {
		if v2, ok := m2[k]; !ok || v1 != v2 {
			return false
		}
	}
	return true
}
