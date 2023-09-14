package slice

import (
	"github.com/matt-repository/matt_golib/common"
	"sort"
)

// Exists ...
func Exists[S ~[]E, E any](arr S, f common.SliceExist[S, E]) bool {
	for i := 0; i < len(arr); i++ {
		if f(arr, i) {
			return true
		}
	}
	return false
}

// Filter ...
func Filter[S ~[]E, E any](arr S, f common.SliceExist[S, E], count int) S {
	ret := make(S, 0)
	for i := 0; i < len(arr); i++ {
		if count <= 0 {
			break
		}
		if f(arr, i) {
			count--
			ret = append(ret, (arr)[i])
		}
	}
	return ret
}

// Select ...
func Select[S1 ~[]E1, S2 []E2, E1 any, E2 any](arr S1, f common.SliceGetElement[S1, E1, E2]) S2 {
	ret := make(S2, 0)
	for i := 0; i < len(arr); i++ {
		v := f(arr, i)
		ret = append(ret, v)
	}
	return ret
}

// ToMap ...
func ToMap[M map[K]V, K comparable, S ~[]E, V, E any](arr S, f common.GetKeyValue[K, V, E]) M {
	ret := make(M)
	for i := 0; i < len(arr); i++ {
		key, value := f(arr[i])
		ret[key] = value
	}
	return ret
}

// Equal ...
func Equal[S ~[]E, E any](arr1, arr2 S, less common.Less[S, E], equal common.EqualElement[E]) bool {
	if len(arr1) != len(arr2) {
		return false
	}
	less1 := func(i, j int) bool {
		return less(arr1, i, j)
	}
	less2 := func(i, j int) bool {
		return less(arr2, i, j)
	}
	sort.Slice(arr1, less1)
	sort.Slice(arr2, less2)

	for i := 0; i < len(arr1); i++ {
		if !equal(arr1[i], arr2[i]) {
			return false
		}
	}
	return true
}

// Differ ...
func Differ[S ~[]E, E any](arr1, arr2 S, equal common.EqualElement[E]) S {
	ret := make(S, 0)
	for _, v1 := range arr1 {
		isExist := Exists(arr2, func(s S, i int) bool {
			return equal(s[i], v1)
		})
		if !isExist {
			ret = append(ret, v1)
		}
	}

	return ret
}

// Intersect ...
func Intersect[S ~[]E, E any](arr1, arr2 S, equal common.EqualElement[E]) S {
	ret := make(S, 0)
	for _, v1 := range arr1 {
		isExist := Exists(arr2, func(s S, i int) bool {
			return equal(s[i], v1)
		})
		if isExist {
			ret = append(ret, v1)
		}
	}
	return ret
}

// Union ...
func Union[S ~[]E, E any](arr1, arr2 S, equal common.EqualElement[E]) S {
	ret := make(S, 0)
	for _, v := range arr1 {
		ret = append(ret, v)
	}
	for _, v2 := range arr2 {
		isExist := Exists(arr1, func(s S, i int) bool {
			return equal(s[i], v2)
		})
		if !isExist {
			ret = append(ret, v2)
		}
	}

	return ret
}

// Distinct ...
func Distinct[S ~[]E, E any](arr S, equal common.EqualElement[E]) S {
	ret := make(S, 0)
	for _, e := range arr {
		isExist := Exists(ret, func(s S, i int) bool {
			return equal(s[i], e)
		})
		if !isExist {
			ret = append(ret, e)
		}
	}
	return ret
}
