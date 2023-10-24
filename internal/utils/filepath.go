package utils

import "strings"

func SanitizeFilePath(path string) string {
	replacements := [][2]string{
		{"'", ""},
		{"...", "."},
		{"..", "."},
		{"$", "."},
		{"!", "."},
		{"?", "."},
		{"…", ""},
	}

	for {
		found := false
		for _, row := range replacements {
			old := row[0]
			new := row[1]
			if strings.Contains(path, old) {
				path = strings.ReplaceAll(path, old, new)
				found = true
			}
		}
		if !found {
			break
		}
	}
	return path
}
