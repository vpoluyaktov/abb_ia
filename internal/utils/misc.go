package utils

import (
)

// Check if a map contains a given key
func HasKey(m map[string]interface{}, key string) bool {
	_, ok := m[key]
	return ok
}

// Checks if a string is present in a slice
func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

// Find an index of elements in a string slice
func GetIndex(s []string, str string) int {
	for i, e := range s {
		 if e == str {  // check the condition if its true return index
				return i
		 }
	}
	return -1
}


