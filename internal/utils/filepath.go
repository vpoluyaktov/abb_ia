package utils

import "strings"

func SanitizeFilePath(path string) string {
	replacements := [][2]string{
		{"'", ""},
		{"..", "."},
		{"...", "."},
		{"$", "."},
		{"!", "."},
		{"?", "."},
	}
	for _, pair := range replacements {
		path = strings.ReplaceAll(path, pair[0], pair[1])
	}
	return path
}
