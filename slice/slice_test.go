package slice

import (
	"fmt"
	"github.com/matt-repository/matt_golib/common"
	_map "github.com/matt-repository/matt_golib/map"
	"strconv"
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

func TestExists(t *testing.T) {
	type args[S []E, E any] struct {
		f common.SliceExist[S, E]
	}
	type testCase[S []E, E any] struct {
		name string
		s    S
		args args[S, E]
		want bool
	}

	//case1: []int
	case1 := testCase[[]int, int]{

		name: "int",
		s:    argsInts,
		args: args[[]int, int]{
			f: func(data []int, i int) bool {
				return data[i] == 1
			},
		},
		want: true,
	}
	if got := Exists(case1.s, case1.args.f); got != case1.want {
		t.Errorf("Exists() = %v, want %v", got, case1.want)
	}

	//case2: []string
	case2 := testCase[[]string, string]{
		name: "string",
		s:    argsStrs,
		args: args[[]string, string]{
			f: func(data []string, i int) bool {
				return data[i] == "1"
			},
		},
		want: true,
	}
	if got := Exists(case2.s, case2.args.f); got != case2.want {
		t.Errorf("Exists() = %v, want %v", got, case2.want)
	}

	//case3: []struct
	case3 := testCase[[]StructA, StructA]{
		name: "Struct",
		s: []StructA{
			{A: 1, B: 1},
		},
		args: args[[]StructA, StructA]{
			f: func(data []StructA, i int) bool {
				return data[i].A == 1 && data[i].B == 1
			},
		},
		want: true,
	}
	if got := Exists(case3.s, case3.args.f); got != case3.want {
		t.Errorf("Exists() = %v, want %v", got, case3.want)
	}

	//case4 []*struct
	case4 := testCase[[]*StructA, *StructA]{
		name: "*StructA",
		s: []*StructA{
			{A: 1, B: 1},
		},
		args: args[[]*StructA, *StructA]{
			f: func(data []*StructA, i int) bool {
				return data[i].A == 1 && data[i].B == 1
			},
		},
		want: true,
	}
	if got := Exists(case4.s, case4.args.f); got != case4.want {
		t.Errorf("Exists() = %v, want %v", got, case4.want)
	}

}

func TestFilter(t *testing.T) {
	type args[S []E, E any] struct {
		f common.SliceExist[S, E]
	}
	type testCase[S []E, E any] struct {
		name string
		s    S
		args args[S, E]
		want S
	}
	//case1: []int
	case1 := testCase[[]int, int]{
		name: "int",
		s:    argsInts,
		args: args[[]int, int]{
			f: func(data []int, i int) bool {
				return data[i] == 1
			},
		},
		want: []int{1},
	}
	got1 := Filter(case1.s, case1.args.f, len(case1.s))
	less1 := func(data []int, i, j int) bool { return data[i] < data[j] }
	equal1 := func(v1, v2 int) bool { return v1 == v2 }
	if !Equal(got1, case1.want, less1, equal1) {
		t.Errorf("Filter() = %v, want %v", got1, case1.want)
	}

	//case2: []string
	case2 := testCase[[]string, string]{
		name: "string",
		s:    argsStrs,
		args: args[[]string, string]{
			f: func(data []string, i int) bool {
				return data[i] == "1"
			},
		},
		want: []string{"1"},
	}
	got2 := Filter(case2.s, case2.args.f, len(case2.s))
	less2 := func(data []string, i, j int) bool { return data[i] < data[j] }
	equal2 := func(v1, v2 string) bool { return v1 == v2 }
	if !Equal(got2, case2.want, less2, equal2) {
		t.Errorf("Filter() = %v, want %v", got2, case2.want)
	}

	//case3: []struct
	case3 := testCase[[]StructA, StructA]{
		name: "struct",
		s:    []StructA{{A: 1, B: 1}, {A: 2, B: 2}},
		args: args[[]StructA, StructA]{
			f: func(data []StructA, i int) bool {
				return data[i].A == 1 && data[i].B == 1
			},
		},
		want: []StructA{{A: 1, B: 1}},
	}
	if got := Filter(case3.s, case3.args.f, len(case3.s)); !Equal(got, case3.want, func(data []StructA, i, j int) bool {
		return fmt.Sprintf("%d:%d", data[i].A, data[i].B) < fmt.Sprintf("%d:%d", data[j].A, data[j].B)
	}, func(v1, v2 StructA) bool { return v1.A == v2.A && v1.B == v2.B }) {
		t.Errorf("Filter() = %v, want %v", got, case3.want)
	}

	//case4: []*struct
	case4 := testCase[[]*StructA, *StructA]{
		name: "*struct",
		s:    []*StructA{{A: 1, B: 1}, {A: 2, B: 2}},
		args: args[[]*StructA, *StructA]{
			f: func(data []*StructA, i int) bool {
				return data[i].A == 1 && data[i].B == 1
			},
		},
		want: []*StructA{{A: 1, B: 1}},
	}
	if got := Filter(case4.s, case4.args.f, len(case4.s)); !Equal(got, case4.want, func(data []*StructA, i, j int) bool {
		return fmt.Sprintf("%d:%d", data[i].A, data[i].B) < fmt.Sprintf("%d:%d", data[j].A, data[j].B)
	}, func(v1, v2 *StructA) bool { return v1.A == v2.A && v1.B == v2.B }) {
		t.Errorf("Filter() = %v, want %v", got, case4.want)
	}

	//case5: []int getCount
	case5 := testCase[[]int, int]{
		name: "int getCount",
		s:    argsInts,
		args: args[[]int, int]{
			f: func(data []int, i int) bool {
				return data[i] >= 1
			},
		},
		want: []int{1},
	}
	got5 := Filter(case5.s, case5.args.f, 1)
	less5 := func(data []int, i, j int) bool { return data[i] < data[j] }
	equal5 := func(v1, v2 int) bool { return v1 == v2 }
	if !Equal(got5, case5.want, less5, equal5) {
		t.Errorf("Filter() = %v, want %v", got5, case5.want)
	}

}

func TestToMap(t *testing.T) {
	type args[E any, K comparable, V any] struct {
		arr []E
		f   common.GetKeyValue[K, V, E]
	}
	type testCase[E any, K comparable, V any, M interface{ ~map[K]V }] struct {
		name string
		args args[E, K, V]
		want M
	}
	//case1: []int=>map[int]int
	case1 := testCase[int, int, int, map[int]int]{
		name: "[]int=>map[int]int",
		args: args[int, int, int]{
			arr: argsInts,
			f: func(v int) (key int, value int) {
				return v, v
			},
		},
		want: argsIntMaps,
	}

	if got := ToMap(case1.args.arr, case1.args.f); !_map.Equal(got, case1.want, func(v1, v2 int) bool {
		return v1 == v2
	}) {
		t.Errorf("Filter() = %v, want %v", got, case1.want)
	}

	//case2 []string=>map[string]string
	case2 := testCase[string, string, string, map[string]string]{
		name: "[]string=>map[string]string",
		args: args[string, string, string]{
			arr: argsStrs,
			f: func(v string) (key string, value string) {
				return v, v
			},
		},
		want: argsStrMaps,
	}

	if got := ToMap(case2.args.arr, case2.args.f); !_map.Equal(got, case2.want, func(v1, v2 string) bool {
		return v1 == v2
	}) {
		t.Errorf("Filter() = %v, want %v", got, case2.want)
	}
}

func TestDiffer(t *testing.T) {
	type args[S []E, E any] struct {
		arr1 S
		arr2 S
		f    common.EqualElement[E]
	}
	type testCase[S []E, E any] struct {
		name string
		args args[S, E]
		want S
	}

	// int
	case1 := testCase[[]int, int]{
		name: "int",
		args: args[[]int, int]{
			arr1: []int{1, 2, 3, 4, 5},
			arr2: []int{1, 3, 4, 5, 6},
			f: func(v1, v2 int) bool {
				return v1 == v2
			},
		},
		want: []int{2},
	}
	got1 := Differ(case1.args.arr1, case1.args.arr2, case1.args.f)

	result1 := Equal(got1, case1.want, func(data []int, i, j int) bool {
		return data[i] < data[j]
	}, func(v1, v2 int) bool {
		return v1 == v2
	})
	if !result1 {
		t.Errorf("case:%s,got:%v,want:%v", case1.name, got1, case1.want)
	}

	// string
	case2 := testCase[[]string, string]{
		name: "string",
		args: args[[]string, string]{
			arr1: []string{"1", "2", "3", "4", "5"},
			arr2: []string{"1", "3", "4", "5", "6"},
			f: func(v1, v2 string) bool {
				return v1 == v2
			},
		},
		want: []string{"2"},
	}
	got2 := Differ(case2.args.arr1, case2.args.arr2, case2.args.f)

	result2 := Equal(got2, case2.want, func(data []string, i, j int) bool {
		return data[i] < data[j]
	}, func(v1, v2 string) bool {
		return v1 == v2
	})
	if !result2 {
		t.Errorf("case:%s,got:%v,want:%v", case2.name, got2, case2.want)
	}

	// Struct
	case3 := testCase[[]StructA, StructA]{
		name: "Struct",
		args: args[[]StructA, StructA]{
			arr1: []StructA{{A: 1, B: 1}, {A: 2, B: 2}, {A: 3, B: 3}, {A: 4, B: 4}, {A: 5, B: 5}},
			arr2: []StructA{{A: 1, B: 1}, {A: 3, B: 3}, {A: 4, B: 4}, {A: 5, B: 5}, {A: 6, B: 6}},
			f: func(v1, v2 StructA) bool {
				return v1 == v2
			},
		},
		want: []StructA{{A: 2, B: 2}},
	}
	got3 := Differ(case3.args.arr1, case3.args.arr2, case3.args.f)

	result3 := Equal(got3, case3.want, func(data []StructA, i, j int) bool {
		return data[i].A < data[j].A && data[i].B < data[i].B
	}, func(v1, v2 StructA) bool {
		return v1 == v2
	})
	if !result3 {
		t.Errorf("case:%s,got:%v,want:%v", case3.name, got3, case3.want)
	}

	// *Struct
	case4 := testCase[[]*StructA, *StructA]{
		name: "*Struct",
		args: args[[]*StructA, *StructA]{
			arr1: []*StructA{{A: 1, B: 1}, {A: 2, B: 2}, {A: 3, B: 3}, {A: 4, B: 4}, {A: 5, B: 5}},
			arr2: []*StructA{{A: 1, B: 1}, {A: 3, B: 3}, {A: 4, B: 4}, {A: 5, B: 5}, {A: 6, B: 6}},
			f: func(v1, v2 *StructA) bool {
				return *v1 == *v2
			},
		},
		want: []*StructA{{A: 2, B: 2}},
	}
	got4 := Differ(case4.args.arr1, case4.args.arr2, case4.args.f)

	result4 := Equal(got4, case4.want, func(data []*StructA, i, j int) bool {
		return data[i].A < data[j].A && data[i].B < data[i].B
	}, func(v1, v2 *StructA) bool {
		return *v1 == *v2
	})
	if !result4 {
		t.Errorf("case:%s,got:%v,want:%v", case4.name, got4, case4.want)
	}
}

func TestInter(t *testing.T) {
	type args[S []E, E any] struct {
		arr1  S
		arr2  S
		equal common.EqualElement[E]
	}
	type testCase[S []E, E any] struct {
		name string
		args args[S, E]
		want S
	}

	// int
	case1 := testCase[[]int, int]{
		name: "int",
		args: args[[]int, int]{
			arr1: []int{1, 2, 3, 4, 5},
			arr2: []int{1, 3, 4, 5, 6},
			equal: func(v1, v2 int) bool {
				return v1 == v2
			},
		},
		want: []int{1, 3, 4, 5},
	}
	got1 := Intersect(case1.args.arr1, case1.args.arr2, case1.args.equal)

	result1 := Equal(got1, case1.want, func(data []int, i, j int) bool {
		return data[i] < data[j]
	}, func(v1, v2 int) bool {
		return v1 == v2
	})
	if !result1 {
		t.Errorf("case:%s,got:%v,want:%v", case1.name, got1, case1.want)
	}

	// string
	case2 := testCase[[]string, string]{
		name: "string",
		args: args[[]string, string]{
			arr1: []string{"1", "2", "3", "4", "5"},
			arr2: []string{"1", "3", "4", "5", "6"},
			equal: func(v1, v2 string) bool {
				return v1 == v2
			},
		},
		want: []string{"1", "3", "4", "5"},
	}
	got2 := Intersect(case2.args.arr1, case2.args.arr2, case2.args.equal)

	result2 := Equal(got2, case2.want, func(data []string, i, j int) bool {
		return data[i] < data[j]
	}, func(v1, v2 string) bool {
		return v1 == v2
	})
	if !result2 {
		t.Errorf("case:%s,got:%v,want:%v", case2.name, got2, case2.want)
	}

	// Struct
	case3 := testCase[[]StructA, StructA]{
		name: "Struct",
		args: args[[]StructA, StructA]{
			arr1: []StructA{{A: 1, B: 1}, {A: 2, B: 2}, {A: 3, B: 3}, {A: 4, B: 4}, {A: 5, B: 5}},
			arr2: []StructA{{A: 1, B: 1}, {A: 3, B: 3}, {A: 4, B: 4}, {A: 5, B: 5}, {A: 6, B: 6}},
			equal: func(v1, v2 StructA) bool {
				return v1 == v2
			},
		},
		want: []StructA{{A: 1, B: 1}, {A: 3, B: 3}, {A: 4, B: 4}, {A: 5, B: 5}},
	}
	got3 := Intersect(case3.args.arr1, case3.args.arr2, case3.args.equal)

	result3 := Equal(got3, case3.want, func(data []StructA, i, j int) bool {
		return data[i].A < data[j].A && data[i].B < data[i].B
	}, func(v1, v2 StructA) bool {
		return v1 == v2
	})
	if !result3 {
		t.Errorf("case:%s,got:%v,want:%v", case3.name, got3, case3.want)
	}

	// *Struct
	case4 := testCase[[]*StructA, *StructA]{
		name: "*Struct",
		args: args[[]*StructA, *StructA]{
			arr1: []*StructA{{A: 1, B: 1}, {A: 2, B: 2}, {A: 3, B: 3}, {A: 4, B: 4}, {A: 5, B: 5}},
			arr2: []*StructA{{A: 1, B: 1}, {A: 3, B: 3}, {A: 4, B: 4}, {A: 5, B: 5}, {A: 6, B: 6}},
			equal: func(v1, v2 *StructA) bool {
				return *v1 == *v2
			},
		},
		want: []*StructA{{A: 1, B: 1}, {A: 3, B: 3}, {A: 4, B: 4}, {A: 5, B: 5}},
	}
	got4 := Intersect(case4.args.arr1, case4.args.arr2, case4.args.equal)

	result4 := Equal(got4, case4.want, func(data []*StructA, i, j int) bool {
		return data[i].A < data[j].A && data[i].B < data[i].B
	}, func(v1, v2 *StructA) bool {
		return *v1 == *v2
	})
	if !result4 {
		t.Errorf("case:%s,got:%v,want:%v", case4.name, got4, case4.want)
	}
}

func TestUnion(t *testing.T) {
	type args[S []E, E any] struct {
		arr1 S
		arr2 S
		f    common.EqualElement[E]
	}
	type testCase[S []E, E any] struct {
		name string
		args args[S, E]
		want S
	}

	// int
	case1 := testCase[[]int, int]{
		name: "int",
		args: args[[]int, int]{
			arr1: []int{1, 2, 3, 4, 5},
			arr2: []int{1, 3, 4, 5, 6},
			f: func(v1, v2 int) bool {
				return v1 == v2
			},
		},
		want: []int{1, 2, 3, 4, 5, 6},
	}
	got1 := Union(case1.args.arr1, case1.args.arr2, case1.args.f)

	result1 := Equal(got1, case1.want, func(data []int, i, j int) bool {
		return data[i] < data[j]
	}, func(v1, v2 int) bool {
		return v1 == v2
	})
	if !result1 {
		t.Errorf("case:%s,got:%v,want:%v", case1.name, got1, case1.want)
	}

	// string
	case2 := testCase[[]string, string]{
		name: "string",
		args: args[[]string, string]{
			arr1: []string{"1", "2", "3", "4", "5"},
			arr2: []string{"1", "3", "4", "5", "6"},
			f: func(v1, v2 string) bool {
				return v1 == v2
			},
		},
		want: []string{"1", "2", "3", "4", "5", "6"},
	}
	got2 := Union(case2.args.arr1, case2.args.arr2, case2.args.f)

	result2 := Equal(got2, case2.want, func(data []string, i, j int) bool {
		return data[i] < data[j]
	}, func(v1, v2 string) bool {
		return v1 == v2
	})
	if !result2 {
		t.Errorf("case:%s,got:%v,want:%v", case2.name, got2, case2.want)
	}

	// Struct
	case3 := testCase[[]StructA, StructA]{
		name: "Struct",
		args: args[[]StructA, StructA]{
			arr1: []StructA{{A: 1, B: 1}, {A: 2, B: 2}, {A: 3, B: 3}, {A: 4, B: 4}, {A: 5, B: 5}},
			arr2: []StructA{{A: 1, B: 1}, {A: 3, B: 3}, {A: 4, B: 4}, {A: 5, B: 5}, {A: 6, B: 6}},
			f: func(v1, v2 StructA) bool {
				return v1 == v2
			},
		},
		want: []StructA{{A: 1, B: 1}, {A: 2, B: 2}, {A: 3, B: 3}, {A: 4, B: 4}, {A: 5, B: 5}, {A: 6, B: 6}},
	}
	got3 := Union(case3.args.arr1, case3.args.arr2, case3.args.f)

	result3 := Equal(got3, case3.want, func(data []StructA, i, j int) bool {
		return data[i].A < data[j].A && data[i].B < data[i].B
	}, func(v1, v2 StructA) bool {
		return v1 == v2
	})
	if !result3 {
		t.Errorf("case:%s,got:%v,want:%v", case3.name, got3, case3.want)
	}

	// *Struct
	case4 := testCase[[]*StructA, *StructA]{
		name: "*Struct",
		args: args[[]*StructA, *StructA]{
			arr1: []*StructA{{A: 1, B: 1}, {A: 2, B: 2}, {A: 3, B: 3}, {A: 4, B: 4}, {A: 5, B: 5}},
			arr2: []*StructA{{A: 1, B: 1}, {A: 3, B: 3}, {A: 4, B: 4}, {A: 5, B: 5}, {A: 6, B: 6}},
			f: func(v1, v2 *StructA) bool {
				return *v1 == *v2
			},
		},
		want: []*StructA{{A: 1, B: 1}, {A: 2, B: 2}, {A: 3, B: 3}, {A: 4, B: 4}, {A: 5, B: 5}, {A: 6, B: 6}},
	}
	got4 := Union(case4.args.arr1, case4.args.arr2, case4.args.f)

	result4 := Equal(got4, case4.want, func(data []*StructA, i, j int) bool {
		return data[i].A < data[j].A && data[i].B < data[i].B
	}, func(v1, v2 *StructA) bool {
		return *v1 == *v2
	})
	if !result4 {
		t.Errorf("case:%s,got:%v,want:%v", case4.name, got4, case4.want)
	}
}

func TestDistinct(t *testing.T) {
	type args[S []E, E any] struct {
		arr   S
		equal common.EqualElement[E]
	}
	type testCase[S []E, E any] struct {
		name string
		args args[S, E]
		want S
	}

	// int
	case1 := testCase[[]int, int]{
		name: "int",
		args: args[[]int, int]{
			arr: []int{1, 2, 2, 4, 4},
			equal: func(v1, v2 int) bool {
				return v1 == v2
			},
		},
		want: []int{1, 2, 4},
	}
	got1 := Distinct(case1.args.arr, case1.args.equal)

	result1 := Equal(got1, case1.want, func(data []int, i, j int) bool {
		return data[i] < data[j]
	}, func(v1, v2 int) bool {
		return v1 == v2
	})
	if !result1 {
		t.Errorf("case:%s,got:%v,want:%v", case1.name, got1, case1.want)
	}

	// string
	case2 := testCase[[]string, string]{
		name: "string",
		args: args[[]string, string]{
			arr: []string{"1", "2", "2", "4", "4"},
			equal: func(v1, v2 string) bool {
				return v1 == v2
			},
		},
		want: []string{"1", "2", "4"},
	}
	got2 := Distinct(case2.args.arr, case2.args.equal)

	result2 := Equal(got2, case2.want, func(data []string, i, j int) bool {
		return data[i] < data[j]
	}, func(v1, v2 string) bool {
		return v1 == v2
	})
	if !result2 {
		t.Errorf("case:%s,got:%v,want:%v", case2.name, got2, case2.want)
	}

	// Struct
	case3 := testCase[[]StructA, StructA]{
		name: "Struct",
		args: args[[]StructA, StructA]{
			arr: []StructA{{A: 1, B: 1}, {A: 2, B: 2}, {A: 2, B: 2}, {A: 4, B: 4}, {A: 4, B: 4}},
			equal: func(v1, v2 StructA) bool {
				return v1 == v2
			},
		},
		want: []StructA{{A: 1, B: 1}, {A: 2, B: 2}, {A: 4, B: 4}},
	}
	got3 := Distinct(case3.args.arr, case3.args.equal)

	result3 := Equal(got3, case3.want, func(data []StructA, i, j int) bool {
		return data[i].A < data[j].A && data[i].B < data[j].B
	}, func(v1, v2 StructA) bool {
		return v1 == v2
	})
	if !result3 {
		t.Errorf("case:%s,got:%v,want:%v", case2.name, got2, case2.want)
	}

	// *Struct
	case4 := testCase[[]*StructA, *StructA]{
		name: "*Struct",
		args: args[[]*StructA, *StructA]{
			arr: []*StructA{{A: 1, B: 1}, {A: 2, B: 2}, {A: 2, B: 2}, {A: 4, B: 4}, {A: 4, B: 4}},
			equal: func(v1, v2 *StructA) bool {
				return v1.A == v2.A && v1.B == v2.B
			},
		},
		want: []*StructA{{A: 1, B: 1}, {A: 2, B: 2}, {A: 4, B: 4}},
	}
	got4 := Distinct(case4.args.arr, case4.args.equal)

	result4 := Equal(got4, case4.want, func(data []*StructA, i, j int) bool {
		return data[i].A < data[j].A && data[i].B < data[j].B
	}, func(v1, v2 *StructA) bool {
		return v1.A == v2.A && v1.B == v2.B
	})
	if !result4 {
		t.Errorf("case:%s,got:%v,want:%v", case4.name, got4, case4.want)
	}

}

func TestSelect(t *testing.T) {
	type args[S1 ~[]E1, E1 any, E2 any] struct {
		arr S1
		f   common.SliceGetElement[S1, E1, E2]
	}
	type testCase[S1 ~[]E1, E1 any, E2 any, S2 []E2] struct {
		name string
		args args[S1, E1, E2]
		want S2
	}

	// int =>string
	case1 := testCase[[]int, int, string, []string]{
		name: "int=>string",
		args: args[[]int, int, string]{
			arr: argsInts,
			f: func(data []int, i int) string {
				return strconv.Itoa(data[i])
			},
		},
		want: argsStrs,
	}
	got1 := Select(case1.args.arr, case1.args.f)
	less1 := func(data []string, i, j int) bool { return data[i] < data[j] }
	equal1 := func(v1, v2 string) bool { return v1 == v2 }
	if !Equal(got1, case1.want, less1, equal1) {
		t.Errorf("Select() = %v, want %v", got1, case1.want)
	}

	// string =>int
	case2 := testCase[[]string, string, int, []int]{
		name: "string =>int",
		args: args[[]string, string, int]{
			arr: argsStrs,
			f: func(data []string, i int) int {
				v, _ := strconv.Atoi(data[i])
				return v
			},
		},
		want: argsInts,
	}
	got2 := Select(case2.args.arr, case2.args.f)
	less2 := func(data []int, i, j int) bool { return data[i] < data[j] }
	equal2 := func(v1, v2 int) bool { return v1 == v2 }
	if !Equal(got2, case2.want, less2, equal2) {
		t.Errorf("Select() = %v, want %v", got2, case2.want)
	}

	// struct =>*struct
	case3 := testCase[[]StructA, StructA, *StructA, []*StructA]{
		name: "struct =>*struct",
		args: args[[]StructA, StructA, *StructA]{
			arr: []StructA{{A: 1, B: 1}, {A: 2, B: 2}},
			f: func(data []StructA, i int) *StructA {
				v := &StructA{A: data[i].A, B: data[i].B}
				return v
			},
		},
		want: []*StructA{{A: 1, B: 1}, {A: 2, B: 2}},
	}
	got3 := Select(case3.args.arr, case3.args.f)
	less3 := func(data []*StructA, i, j int) bool { return data[i].A < data[j].A && data[i].B < data[i].B }
	equal3 := func(v1, v2 *StructA) bool { return *v1 == *v2 }
	if !Equal(got3, case3.want, less3, equal3) {
		t.Errorf("Select() = %v, want %v", got3, case3.want)
	}

}
