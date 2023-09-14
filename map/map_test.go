package _map

import (
	"fmt"
	"github.com/matt-repository/matt_golib/common"
	"github.com/matt-repository/matt_golib/slice"
	"testing"
)

type StructA struct {
	A1 int
	A2 string
}

func TestToSlice(t *testing.T) {
	//case1: [int]string => []string
	m1 := map[int]string{
		1: "二",
		2: "一",
	}
	slice1 := ToSlice(m1, func(k int, v string) string {
		return fmt.Sprintf("%d:%s", k, v)
	})
	want1 := []string{"1:二", "2:一"}

	equalResult1 := slice.Equal(slice1, want1, func(data []string, i, j int) bool { return data[i] < data[j] }, func(v1, v2 string) bool { return v1 == v2 })
	if !equalResult1 {
		t.Errorf("MapToSlice() = %v, want1 %v", slice1, want1)
	}

	//case2: [string]int => []string
	m2 := map[string]int{
		"二": 1,
		"一": 2,
	}
	slice2 := ToSlice(m2, func(k string, v int) string {
		return fmt.Sprintf("%s:%d", k, v)
	})
	want2 := []string{"二:1", "一:2"}

	equalResult2 := slice.Equal(slice2, want2, func(data []string, i, j int) bool { return data[i] < data[j] }, func(v1, v2 string) bool { return v1 == v2 })
	if !equalResult2 {
		t.Errorf("MapToSlice() = %v, want2 %v", slice2, want2)
	}

	//case3: [string]StructA => []string
	m3 := map[string]StructA{
		"二": {
			A1: 2,
			A2: "2",
		},
		"一": {
			A1: 1,
			A2: "1",
		},
	}
	slice3 := ToSlice(m3, func(k string, v StructA) string {
		return fmt.Sprintf("%s:[%d,%s]", k, v.A1, v.A2)
	})
	want3 := []string{"二:[2,2]", "一:[1,1]"}

	equalResult3 := slice.Equal(slice3, want3, func(data []string, i, j int) bool { return data[i] < data[j] }, func(v1, v2 string) bool { return v1 == v2 })
	if !equalResult3 {
		t.Errorf("MapToSlice() = %v, want3 %v", slice3, want3)
	}

	//case4: [int]string=>[]StructA
	m4 := map[int]string{
		1: "二",
		2: "一",
	}
	slice4 := ToSlice(m4, func(k int, v string) StructA {
		return StructA{
			A1: k,
			A2: v,
		}
	})
	want4 := []StructA{{
		A1: 1,
		A2: "二",
	}, {
		A1: 2,
		A2: "一",
	}}

	equalResult4 := slice.Equal(slice4, want4, func(data []StructA, i, j int) bool {
		return fmt.Sprintf("%d:%s", data[i].A1, data[i].A2) < fmt.Sprintf("%d:%s", data[j].A1, data[j].A2)
	}, func(v1, v2 StructA) bool { return v1.A1 == v2.A1 && v1.A2 == v2.A2 })
	if !equalResult4 {
		t.Errorf("MapToSlice() = %v, want4 %v", slice4, want4)
	}

}

func TestEqual(t *testing.T) {
	type args[M ~map[K]V, K comparable, V any] struct {
		m1    M
		m2    M
		equal common.EqualElement[V]
	}
	type testCase[M ~map[K]V, K comparable, V any] struct {
		name string
		args args[M, K, V]
		want bool
	}

	//case1 :int:int
	case1 := testCase[map[int]int, int, int]{
		name: "int:int",
		args: args[map[int]int, int, int]{
			m1: map[int]int{1: 1, 2: 2, 3: 3},
			m2: map[int]int{1: 1, 2: 2, 3: 3},
			equal: func(v1, v2 int) bool {
				return v1 == v2
			},
		},
		want: true,
	}
	got1 := Equal(case1.args.m1, case1.args.m2, case1.args.equal)
	if got1 != case1.want {
		t.Errorf("Equal() = %v, want1 %v", got1, case1.want)
	}

	//case2 :string:string
	case2 := testCase[map[string]string, string, string]{
		name: "string:string",
		args: args[map[string]string, string, string]{
			m1: map[string]string{"1": "1", "2": "2", "3": "3"},
			m2: map[string]string{"1": "1", "2": "2", "3": "3"},
			equal: func(v1, v2 string) bool {
				return v1 == v2
			},
		},
		want: true,
	}
	got2 := Equal(case2.args.m1, case2.args.m2, case2.args.equal)
	if got2 != case2.want {
		t.Errorf("Equal() = %v, want2 %v", got2, case2.want)
	}

	//case3 :int:string
	case3 := testCase[map[int]string, int, string]{
		name: "int:string",
		args: args[map[int]string, int, string]{
			m1: map[int]string{1: "1", 2: "2", 3: "3"},
			m2: map[int]string{1: "1", 2: "2", 3: "3"},
			equal: func(v1, v2 string) bool {
				return v1 == v2
			},
		},
		want: true,
	}
	got3 := Equal(case3.args.m1, case3.args.m2, case3.args.equal)
	if got3 != case3.want {
		t.Errorf("Equal() = %v, want3 %v", got3, case3.want)
	}

	//case4 :int:StructA
	case4 := testCase[map[int]StructA, int, StructA]{
		name: "int:StructA",
		args: args[map[int]StructA, int, StructA]{
			m1: map[int]StructA{1: {A1: 1, A2: "1"}, 2: {A1: 2, A2: "2"}, 3: {A1: 3, A2: "3"}},
			m2: map[int]StructA{1: {A1: 1, A2: "1"}, 2: {A1: 2, A2: "2"}, 3: {A1: 3, A2: "3"}},
			equal: func(v1, v2 StructA) bool {
				return v1 == v2
			},
		},
		want: true,
	}
	got4 := Equal(case4.args.m1, case4.args.m2, case4.args.equal)
	if got4 != case4.want {
		t.Errorf("Equal() = %v, want4 %v", got4, case4.want)
	}

	//case5 :int:*StructA
	case5 := testCase[map[int]*StructA, int, *StructA]{
		name: "int:*StructA",
		args: args[map[int]*StructA, int, *StructA]{
			m1: map[int]*StructA{1: {A1: 1, A2: "1"}, 2: {A1: 2, A2: "2"}, 3: {A1: 3, A2: "3"}},
			m2: map[int]*StructA{1: {A1: 1, A2: "1"}, 2: {A1: 2, A2: "2"}, 3: {A1: 3, A2: "3"}},
			equal: func(v1, v2 *StructA) bool {
				return v1 == v2
			},
		},
		want: false,
	}
	got5 := Equal(case5.args.m1, case5.args.m2, case5.args.equal)
	if got5 != case5.want {
		t.Errorf("Equal() = %v, want5 %v", got5, case5.want)
	}

}
