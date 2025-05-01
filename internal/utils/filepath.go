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
		{"#", "N"},
		{"[", ""},
		{"]", ""},
		{":", "."},
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

// TODO: Refactor this (regex?)
func SanitizeMp3FileName(fileName string) string {
	replacements := [][2]string{
		{"_64kb", ""},
		{"_24kb", ""},
		{"_32kb", ""},
		{"_40kb", ""},
		{"_48kb", ""},
		{"_56kb", ""},
		{"_64kb", ""},
		{"_80kb", ""},
		{"_96kb", ""},
		{"_112kb", ""},
		{"_128kb", ""},
		{"_144kb", ""},
		{"_160kb", ""},
		{"_224kb", ""},
		{"_256kb", ""},
		{"_320kb", ""},
		{"_vbr", ""},
	}

	for {
		found := false
		for _, row := range replacements {
			old := row[0]
			new := row[1]
			if strings.Contains(fileName, old) {
				fileName = strings.ReplaceAll(fileName, old, new)
			}
		}
		if !found {
			break
		}
	}
	return fileName
}
