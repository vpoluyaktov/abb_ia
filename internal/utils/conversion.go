package utils

import "strconv"

// check if the string is a number (int, float)
func IsNumber(str string) bool {
	_, err := strconv.ParseFloat(str, 64)
	return err == nil
}

// convert string to int and ignore an error
func ToInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

// convert any type of number to a string and ignore an error
func ToString(num interface{}) string {
	switch num := num.(type) {
	case int:
		return strconv.Itoa(num)
	case int32:
		return strconv.FormatInt(int64(num), 10)
	case int64:
		return strconv.FormatInt(num, 10)
	case float32:
		return strconv.FormatFloat(float64(num), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(num, 'f', -1, 64)
	default:
		return ""
	}
}