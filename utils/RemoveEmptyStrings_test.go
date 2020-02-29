package utils

import (
	"testing"
)

func sliceCompare(arr1 []string, arr2 []string) bool {
	if len(arr1) != len(arr2) {
		return false
	}

	for i := range arr1 {
		if arr1[i] != arr2[i] {
			return false
		}
	}

	return true
}

func TestRemoveEmptyStrings(t *testing.T) {
	input := []string{"hello", "", "     world  "}
	output := RemoveEmptyStrings(input)
	expected := []string{"hello", "world"}

	if !sliceCompare(output, expected) {
		t.Errorf("RemoveEmptyStrings, got: %+q, expected: %+q", output, expected)
	}
}
