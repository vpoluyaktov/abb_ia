package utils

import (
	"regexp"
	"strconv"
	"strings"
)

// romanValues maps Roman numeral symbols to their decimal values
var romanValues = map[byte]int{
	'I': 1,
	'V': 5,
	'X': 10,
	'L': 50,
	'C': 100,
	'D': 500,
	'M': 1000,
}

// romanPattern matches Roman numerals at the start of a string
var (
	romanPattern = regexp.MustCompile(`^[IVXLCDM]+`)
	numberPattern = regexp.MustCompile(`^\d+`)
)

// RomanToDecimal converts a Roman numeral string to decimal
func RomanToDecimal(roman string) int {
	if roman == "" {
		return 0
	}

	total := 0
	prevValue := 0

	// Process from right to left
	for i := len(roman) - 1; i >= 0; i-- {
		currValue := romanValues[roman[i]]

		// If current value is greater than or equal to previous value,
		// add it to total (e.g., VI = 5 + 1 = 6)
		if currValue >= prevValue {
			total += currValue
		} else {
			// If current value is less than previous value,
			// subtract it from total (e.g., IV = 5 - 1 = 4)
			total -= currValue
		}

		prevValue = currValue
	}

	return total
}

// ExtractLeadingRoman extracts a leading Roman numeral from a string.
// Returns the decimal value of the Roman numeral and the remaining string.
func ExtractLeadingRoman(s string) (int, string) {
	// Find Roman numeral at start of string
	match := romanPattern.FindString(s)
	if match == "" {
		return 0, s
	}

	// Convert to decimal
	num := RomanToDecimal(match)
	if num == 0 {
		return 0, s
	}

	// Return remaining string after Roman numeral
	rest := s[len(match):]
	return num, rest
}

// ExtractLeadingNumber extracts the first Arabic number found at the start of a string and returns it along with remaining text.
// Returns 0 if no number is found at the start.
func ExtractLeadingNumber(s string) (int, string) {
	// Find first number at start of string
	match := numberPattern.FindString(s)
	if match == "" {
		return 0, s
	}

	// Convert to integer
	num, err := strconv.Atoi(match)
	if err != nil {
		return 0, s
	}

	// Return remaining string after number
	rest := s[len(match):]
	return num, rest
}

// compareArabicOrder compares two strings using natural sort order with Arabic numbers.
// Returns true if s1 should come before s2.
func compareArabicOrder(s1, s2 string) bool {
	// Split strings into words
	words1 := strings.Fields(s1)
	words2 := strings.Fields(s2)

	// Compare word by word
	for i := 0; i < len(words1) && i < len(words2); i++ {
		// If words are identical, continue to next word
		if words1[i] == words2[i] {
			continue
		}

		// Try to extract numbers from both words
		num1, _ := ExtractLeadingNumber(words1[i])
		num2, _ := ExtractLeadingNumber(words2[i])

		// If both words have numbers
		if num1 > 0 && num2 > 0 {
			return num1 < num2
		}

		// If one has a number and the other doesn't, number comes first
		if (num1 == 0) != (num2 == 0) {
			return num1 != 0
		}

		// If no numbers, sort lexicographically
		return words1[i] < words2[i]
	}

	// If all words match up to the shorter string, shorter string comes first
	return len(words1) < len(words2)
}

// compareRomanOrder compares two strings using natural sort order with Roman numerals.
// Returns true if s1 should come before s2.
func compareRomanOrder(s1, s2 string) bool {
	// Split strings into words
	words1 := strings.Fields(s1)
	words2 := strings.Fields(s2)

	// Compare word by word
	for i := 0; i < len(words1) && i < len(words2); i++ {
		// If words are identical, continue to next word
		if words1[i] == words2[i] {
			continue
		}

		// Try to extract numbers from both words
		num1, _ := ExtractLeadingRoman(words1[i])
		num2, _ := ExtractLeadingRoman(words2[i])

		// If both words have numbers
		if num1 > 0 && num2 > 0 {
			return num1 < num2
		}

		// If one has a number and the other doesn't, number comes first
		if (num1 == 0) != (num2 == 0) {
			return num1 != 0
		}

		// If no numbers, sort lexicographically
		return words1[i] < words2[i]
	}

	// If all words match up to the shorter string, shorter string comes first
	return len(words1) < len(words2)
}

// CompareNaturalOrder compares two strings using natural sort order.
// Handles either Arabic numbers (1, 2, 10) or Roman numerals (I, II, X), but not mixed.
// The first string determines whether we're in Roman or Arabic mode.
// Returns true if s1 should come before s2.
func CompareNaturalOrder(s1, s2 string) bool {
	// First, detect if we're in Roman numeral mode by checking the first string
	romNum1, _ := ExtractLeadingRoman(s1)
	isRomanMode := romNum1 > 0

	if isRomanMode {
		return compareRomanOrder(s1, s2)
	}
	return compareArabicOrder(s1, s2)
}
