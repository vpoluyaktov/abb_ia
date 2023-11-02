package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/vpoluyaktov/abb_ia/internal/logger"
)

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
func SecondsToTime(sec float64) string {
	ss := int(sec) % 60
	sec /= 60
	mm := int(sec) % 60
	sec /= 60
	hh := int(sec)
	time := fmt.Sprintf("%2d:%02d:%02d", hh, mm, ss)
	return time
}

// Convert bytes to Gb, Kb, Mb, string format
func BytesToHuman(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d b", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cb", float64(b)/float64(div), "KMGTPE"[exp])
}

func HumanToBytes(size string) (int64, error) {
	// Remove leading/trailing spaces and convert to lowercase
	size = strings.TrimSpace(strings.ToLower(size))

	// Extract number and size suffix by splitting on spaces
	parts := strings.Fields(size)
	if len(parts) != 2 {
		// Handle formats where number and size suffix are concatenated
		numStr, suffix := separateNumberAndSuffix(size)
		parts = []string{numStr, suffix}
	}

	// Extract number and parse as float64
	numStr, suffix := parts[0], parts[1]
	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid file size")
	}

	// Convert number to bytes based on suffix
	var bytes int64
	switch suffix {
	case "b":
		bytes = int64(num)
	case "kb":
		bytes = int64(num * 1024)
	case "mb":
		bytes = int64(num * 1024 * 1024)
	case "gb":
		bytes = int64(num * 1024 * 1024 * 1024)
	case "tb":
		bytes = int64(num * 1024 * 1024 * 1024 * 1024)
	default:
		return 0, fmt.Errorf("invalid size suffix")
	}

	return bytes, nil
}

// Separates the number and size suffix from concatenated format
func separateNumberAndSuffix(size string) (string, string) {
	numStr := ""
	suffix := ""
	for i, ch := range size {
		if ch >= '0' && ch <= '9' || ch == '.' {
			numStr += string(ch)
		} else {
			suffix = strings.TrimSpace(size[i:])
			break
		}
	}
	return numStr, suffix
}

// Convert download speed to Gb/s, Kb/s, Mb/s, string format
func SpeedToHuman(b int64) string {
	bytesH := BytesToHuman(b)
	return fmt.Sprintf("%s%s", bytesH, "/s")
}

func FirstN(s string, n int) string {
	if len(s) > n {
		return s[:n] + "…"
	}
	return s
}
