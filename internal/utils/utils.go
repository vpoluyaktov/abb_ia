package utils

import (
	"fmt"
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

// Convert time in seconds to HH:MM:SS striing format
func SecondsToTime(sec float64) (string, error) {
	ss := int(sec) % 60
	sec /= 60
	mm := int(sec) % 60
	sec /= 60
	hh := int(sec)
	time := fmt.Sprintf("%2d:%02d:%02d", hh, mm, ss)
	return time, nil
}

// Convert bytes to Gb, Kb, Mb, string format
func BytesToHuman(b int64) (string, error) {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b), nil
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp]), nil
}

// Convert download speed to Gb/s, Kb/s, Mb/s, string format
func SpeedToHuman(b int64) (string, error) {
	bytesH, _ := BytesToHuman(b)
	return fmt.Sprintf("%s%s", bytesH, "/s"), nil
}


func FirstN(s string, n int) string {
	if len(s) > n {
		return s[:n] + "…"
	}
	return s
}
