package utils

import "testing"

func TestFilterStringArray(t *testing.T) {
	arr := []string{"a", "b", "c", "d"}
	filterFunc := func(v string) bool {
		return v != "a"
	}

	arr = FilterStringArray(arr, filterFunc)
	if len(arr) != 3 {
		t.Errorf("Result array should filter out letter a, but got %v", arr)
	}
}
