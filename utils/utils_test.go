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

func TestSubtractStringArray(t *testing.T) {
	firstArr := []string{"a", "b", "c", "d"}
	secondArr := []string{"c", "d"}

	result := SubtractStringArray(firstArr, secondArr)

	if len(result) != 2 {
		t.Errorf("Should get array of 2 elements. Got %v", result)
	}

	for i, v := range []string{"a", "b"} {
		if result[i] != v {
			t.Errorf("Expect element at %d is %s, but got %s", i, v, result[i])
		}
	}
}
