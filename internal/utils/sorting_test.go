package utils

import (
	"sort"
	"testing"
)

func TestRomanToDecimal(t *testing.T) {
	tests := []struct {
		name     string
		roman    string
		expected int
	}{
		{"empty string", "", 0},
		{"single digits", "I", 1},
		{"simple additive", "VII", 7},
		{"simple subtractive", "IV", 4},
		{"complex number", "MCMXCIX", 1999},
		{"all letters", "MDCLXVI", 1666},
		{"repeated letters", "XXXIX", 39},
		{"common numbers", "III", 3},
		{"subtractive pairs", "XL", 40},
		{"large number", "MMXXI", 2021},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RomanToDecimal(tt.roman)
			if got != tt.expected {
				t.Errorf("RomanToDecimal(%q) = %d; want %d", tt.roman, got, tt.expected)
			}
		})
	}
}

func TestExtractLeadingRoman(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedNum    int
		expectedRest   string
	}{
		{"empty string", "", 0, ""},
		{"just roman", "IV", 4, ""},
		{"roman with text", "IV chapter", 4, " chapter"},
		{"roman with number", "IV 2", 4, " 2"},
		{"text only", "chapter", 0, "chapter"},
		{"complex roman", "XXXIX stories", 39, " stories"},
		{"mixed case", "iV chapter", 0, "iV chapter"}, // should be case sensitive
		{"partial roman", "IVx", 4, "x"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			num, rest := ExtractLeadingRoman(tt.input)
			if num != tt.expectedNum || rest != tt.expectedRest {
				t.Errorf("ExtractLeadingRoman(%q) = (%d, %q); want (%d, %q)",
					tt.input, num, rest, tt.expectedNum, tt.expectedRest)
			}
		})
	}
}

func TestExtractLeadingNumber(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedNum    int
		expectedRest   string
	}{
		{"empty string", "", 0, ""},
		{"just number", "42", 42, ""},
		{"number with text", "42 chapter", 42, " chapter"},
		{"text with number", "chapter 42", 0, "chapter 42"},
		{"multiple numbers", "42 43", 42, " 43"},
		{"text only", "chapter", 0, "chapter"},
		{"zero", "0", 0, ""},
		{"leading zeros", "007", 7, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			num, rest := ExtractLeadingNumber(tt.input)
			if num != tt.expectedNum || rest != tt.expectedRest {
				t.Errorf("ExtractLeadingNumber(%q) = (%d, %q); want (%d, %q)",
					tt.input, num, rest, tt.expectedNum, tt.expectedRest)
			}
		})
	}
}

func TestCompareArabicOrder(t *testing.T) {
	tests := []struct {
		name string
		s1   string
		s2   string
		want bool
	}{
		{
			name: "simple_numbers",
			s1:   "2",
			s2:   "10",
			want: true,
		},
		{
			name: "text_with_numbers",
			s1:   "chapter 2",
			s2:   "chapter 10",
			want: true,
		},
		{
			name: "equal_numbers",
			s1:   "chapter 2",
			s2:   "chapter 2",
			want: false,
		},
		{
			name: "pure_text",
			s1:   "abc",
			s2:   "def",
			want: true,
		},
		{
			name: "mixed_text_and_numbers",
			s1:   "chapter 2 section",
			s2:   "chapter 10 intro",
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compareArabicOrder(tt.s1, tt.s2); got != tt.want {
				t.Errorf("compareArabicOrder(%q, %q) = %v; want %v", tt.s1, tt.s2, got, tt.want)
			}
		})
	}
}

func TestCompareRomanOrder(t *testing.T) {
	tests := []struct {
		name string
		s1   string
		s2   string
		want bool
	}{
		{
			name: "simple_roman",
			s1:   "II",
			s2:   "X",
			want: true,
		},
		{
			name: "text_with_roman",
			s1:   "chapter II",
			s2:   "chapter X",
			want: true,
		},
		{
			name: "equal_roman",
			s1:   "chapter V",
			s2:   "chapter V",
			want: false,
		},
		{
			name: "complex_roman",
			s1:   "chapter XXXIX",
			s2:   "chapter XL",
			want: true,
		},
		{
			name: "pure_text",
			s1:   "abc",
			s2:   "def",
			want: true,
		},
		{
			name: "mixed_text_and_roman",
			s1:   "chapter II section",
			s2:   "chapter X intro",
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compareRomanOrder(tt.s1, tt.s2); got != tt.want {
				t.Errorf("compareRomanOrder(%q, %q) = %v; want %v", tt.s1, tt.s2, got, tt.want)
			}
		})
	}
}

func TestCompareNaturalOrder(t *testing.T) {
	tests := []struct {
		name string
		s1   string
		s2   string
		want bool
	}{
		{
			name: "arabic_mode",
			s1:   "chapter 2",
			s2:   "chapter 10",
			want: true,
		},
		{
			name: "roman_mode",
			s1:   "chapter II",
			s2:   "chapter X",
			want: true,
		},
		{
			name: "pure_text",
			s1:   "abc",
			s2:   "def",
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CompareNaturalOrder(tt.s1, tt.s2); got != tt.want {
				t.Errorf("CompareNaturalOrder(%q, %q) = %v; want %v", tt.s1, tt.s2, got, tt.want)
			}
		})
	}
}

func TestSortingWithArabicNumbers(t *testing.T) {
	items := []string{
		"Chapter 10",
		"Chapter 2",
		"Chapter 4",
		"Chapter 5",
	}

	expected := []string{
		"Chapter 2",
		"Chapter 4",
		"Chapter 5",
		"Chapter 10",
	}

	sort.SliceStable(items, func(i, j int) bool {
		return compareArabicOrder(items[i], items[j])
	})

	for i := range items {
		if items[i] != expected[i] {
			t.Errorf("\nAt position %d:\n\tgot:  %q\n\twant: %q", i, items[i], expected[i])
		}
	}
}

func TestSortingWithRomanNumerals(t *testing.T) {
	items := []string{
		"Chapter V",
		"Chapter X",
		"Chapter II",
		"Chapter IV",
		"Chapter I",
	}

	expected := []string{
		"Chapter I",
		"Chapter II",
		"Chapter IV",
		"Chapter V",
		"Chapter X",
	}

	sort.SliceStable(items, func(i, j int) bool {
		return CompareNaturalOrder(items[i], items[j])
	})

	for i := range items {
		if items[i] != expected[i] {
			t.Errorf("\nAt position %d:\n\tgot:  %q\n\twant: %q", i, items[i], expected[i])
		}
	}
}
