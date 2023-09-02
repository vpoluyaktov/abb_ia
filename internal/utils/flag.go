package utils

import "flag"

// Check if a command link flag passed
func IsFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
			if f.Name == name {
					found = true
			}
	})
	return found
}