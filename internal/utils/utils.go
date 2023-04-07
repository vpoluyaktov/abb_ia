package utils

import (
	"strconv"
	"strings"

	"github.com/vpoluyaktov/audiobook_creator_IA/internal/logger"
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

// Convert time in HH:MM:SS or SSSSS.MI string format to seconds
func TimeToSeconds(t string) (float64, error) {
	slices := strings.Split(t, ":")
	if len(slices) > 1 { // HH:MM:SS format
		var sec float64 = 0
		for _, s := range slices {
			sDuration, err := strconv.ParseFloat(s, 64)
			if err != nil {
				logger.Error("Can't convert time to seconds: " + t)
				return 0, err
			}
			sec = sec*60 + sDuration
		}
		return sec, nil
	} else { // SSSSS.MI format
		sec, err := strconv.ParseFloat(t, 64)
		if err != nil {
			return 0, err
		}
		return sec, nil
	}
}
