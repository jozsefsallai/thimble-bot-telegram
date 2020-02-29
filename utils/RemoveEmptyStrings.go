package utils

import "strings"

// RemoveEmptyStrings will remove all empty strings from a string slice
func RemoveEmptyStrings(arr []string) []string {
	var output []string

	for _, current := range arr {
		actual := strings.TrimSpace(current)
		if len(actual) > 0 {
			output = append(output, actual)
		}
	}

	return output
}
