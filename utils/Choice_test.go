package utils

import (
	"testing"
)

func isOK(str string) bool {
	switch str {
	case
		"joe",
		"bryan",
		"eli":
		return true
	}

	return false
}

func TestChoiceString(t *testing.T) {
	input := []string{"joe", "bryan", "eli"}
	choice := ChoiceString(input)
	result := isOK(choice)

	if result == false {
		t.Errorf("ChoiceString, got: %t, expected: %t", result, true)
	}
}
