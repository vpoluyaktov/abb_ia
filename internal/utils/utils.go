package utils

// Check if a map contains a given key
func HasKey(m map[string]interface{}, key string) bool {
	_, ok := m[key]
	return ok
}
