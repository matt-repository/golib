package util

import (
	"testing"
)

var (
	argsInts    = []int{1, 2, 3, 4, 5, 6}
	argsIntMaps = map[int]int{
		1: 1,
		2: 2,
		3: 3,
		4: 4,
		5: 5,
		6: 6,
	}
	argsStrs    = []string{"1", "2", "3", "4", "5", "6"}
	argsStrMaps = map[string]string{
		"1": "1",
		"2": "2",
		"3": "3",
		"4": "4",
		"5": "5",
		"6": "6",
	}
)

type StructA struct {
	A, B int
}

func Test_slice_Contains(t *testing.T) {
	type args[T any] struct {
		f condition[T]
	}
	type testCase[T any] struct {
		name string
		s    []T
		args args[T]
		want bool
	}
	case1 := testCase[int]{

		name: "int",
		s:    argsInts,
		args: args[int]{
			f: func(data []int, i int) bool {
				return data[i] == 1
			},
		},
		want: true,
	}
	if got := Contains(case1.s, case1.args.f); got != case1.want {
		t.Errorf("Contains() = %v, want %v", got, case1.want)
	}

	case2 := testCase[string]{

		name: "string",
		s:    argsStrs,
		args: args[string]{
			f: func(data []string, i int) bool {
				return data[i] == "1"
			},
		},
		want: true,
	}
	if got := Contains(case2.s, case2.args.f); got != case2.want {
		t.Errorf("Contains() = %v, want %v", got, case2.want)
	}

	case3 := testCase[StructA]{
		name: "struct",
		s: []StructA{
			{A: 1, B: 1},
		},
		args: args[StructA]{
			f: func(data []StructA, i int) bool {
				return data[i].A == 1 && data[i].B == 1
			},
		},
		want: true,
	}
	if got := Contains(case3.s, case3.args.f); got != case3.want {
		t.Errorf("Contains() = %v, want %v", got, case3.want)
	}

}

func Test_slice_Filter(t *testing.T) {
	type args[T any] struct {
		f condition[T]
	}
	type testCase[T any] struct {
		name string
		s    []T
		args args[T]
		want []T
	}
	case1 := testCase[int]{

		name: "int",
		s:    argsInts,
		args: args[int]{
			f: func(data []int, i int) bool {
				return data[i] == 1
			},
		},
		want: []int{1},
	}
	if got := Filter(case1.s, case1.args.f); !EqualSlice(got, case1.want, func(data []int, i, j int) bool { return data[i] < data[j] }, func(v1, v2 int) bool { return v1 == v2 }) {
		t.Errorf("Filter() = %v, want %v", got, case1.want)
	}

	case2 := testCase[string]{
		name: "string",
		s:    argsStrs,
		args: args[string]{
			f: func(data []string, i int) bool {
				return data[i] == "1"
			},
		},
		want: []string{"1"},
	}
	if got := Filter(case2.s, case2.args.f); !EqualSlice(got, case2.want, func(data []string, i, j int) bool { return data[i] < data[j] }, func(v1, v2 string) bool { return v1 == v2 }) {
		t.Errorf("Filter() = %v, want %v", got, case2.want)
	}
}

func TestSliceToMap(t *testing.T) {
	type args[E any, K comparable, V any] struct {
		arr []E
		f   getKeyValue[K, V, E]
	}
	type testCase[E any, K comparable, V any, M interface{ ~map[K]V }] struct {
		name string
		args args[E, K, V]
		want M
	}

	case1 := testCase[int, int, int, map[int]int]{
		name: "int",
		args: args[int, int, int]{
			arr: argsInts,
			f: func(v int) (key int, value int) {
				return v, v
			},
		},
		want: argsIntMaps,
	}

	if got := SliceToMap[map[int]int, int, int](case1.args.arr, case1.args.f); !EqualMap(got, case1.want) {
		t.Errorf("Filter() = %v, want %v", got, case1.want)
	}

	case2 := testCase[string, string, string, map[string]string]{
		name: "string",
		args: args[string, string, string]{
			arr: argsStrs,
			f: func(v string) (key string, value string) {
				return v, v
			},
		},
		want: argsStrMaps,
	}

	if got := SliceToMap[map[string]string, string, string](case2.args.arr, case2.args.f); !EqualMap(got, case2.want) {
		t.Errorf("Filter() = %v, want %v", got, case2.want)
	}
}
