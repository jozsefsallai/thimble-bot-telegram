package utils

// HasPermission checks if a user has a certain permission
func HasPermission(userID int, permission []int) bool {
	for _, user := range permission {
		if user == userID {
			return true
		}
	}

	return false
}
