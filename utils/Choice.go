package utils

import "math/rand"

// ChoiceString picks a random string from a string array
func ChoiceString(arr []string) string {
	return arr[rand.Intn(len(arr))]
}
