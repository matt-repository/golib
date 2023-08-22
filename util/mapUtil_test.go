package util

import (
	"fmt"
	"testing"
)

func TestMapToSlice(t *testing.T) {
	m1 := map[int]string{
		1: "二",
		2: "一",
	}
	slice := MapToSlice(m1, func(k int, v string) string {
		return fmt.Sprintf("%d:%s", k, v)
	})
	want := []string{"1:二", "2:一"}

	equalResult := EqualSlice(slice, want, func(data []string, i, j int) bool { return data[i] < data[j] }, func(v1, v2 string) bool { return v1 == v2 })
	if !equalResult {
		t.Errorf("MapToSlice() = %v, want %v", slice, want)
	}
}
