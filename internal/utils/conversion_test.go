package utils

import "testing"

func TestIsNumber(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want bool
	}{
		{"empty string", "", false},
		{"integer", "123", true},
		{"negative integer", "-123", true},
		{"float", "123.456", true},
		{"negative float", "-123.456", true},
		{"scientific notation", "1.23e5", true},
		{"text", "abc", false},
		{"mixed", "123abc", false},
		{"special chars", "!@#", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNumber(tt.str); got != tt.want {
				t.Errorf("IsNumber(%q) = %v, want %v", tt.str, got, tt.want)
			}
		})
	}
}

func TestToInt(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want int
	}{
		{"empty string", "", 0},
		{"zero", "0", 0},
		{"positive", "123", 123},
		{"negative", "-123", -123},
		{"invalid", "abc", 0},
		{"float", "123.456", 0},
		{"mixed", "123abc", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToInt(tt.s); got != tt.want {
				t.Errorf("ToInt(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func TestToString(t *testing.T) {
	tests := []struct {
		name string
		num  interface{}
		want string
	}{
		{"int", 123, "123"},
		{"int32", int32(123), "123"},
		{"int64", int64(123), "123"},
		{"float32", float32(123.456), "123.456"},
		{"float64", 123.456, "123.456"},
		{"negative int", -123, "-123"},
		{"negative float", -123.456, "-123.456"},
		{"zero", 0, "0"},
		{"invalid type", "123", ""},
		{"nil", nil, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToString(tt.num); got != tt.want {
				t.Errorf("ToString(%v) = %v, want %v", tt.num, got, tt.want)
			}
		})
	}
}
