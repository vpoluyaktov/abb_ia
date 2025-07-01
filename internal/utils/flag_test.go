package utils

import (
	"flag"
	"os"
	"testing"
)

func TestIsFlagPassed(t *testing.T) {
	// Save original args and flags
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	}()

	tests := []struct {
		name     string
		args     []string
		flagName string
		want     bool
	}{
		{
			name:     "no flags",
			args:     []string{"program"},
			flagName: "test",
			want:     false,
		},
		{
			name:     "flag exists and set",
			args:     []string{"program", "-test=value"},
			flagName: "test",
			want:     true,
		},
		{
			name:     "flag exists but not set",
			args:     []string{"program"},
			flagName: "test",
			want:     false,
		},
		{
			name:     "different flag set",
			args:     []string{"program", "-other=value"},
			flagName: "test",
			want:     false,
		},
		{
			name:     "boolean flag set",
			args:     []string{"program", "-test"},
			flagName: "test",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags for each test
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			if tt.name == "boolean flag set" {
				flag.Bool("test", false, "test flag")
			} else {
				flag.String("test", "", "test flag")
			}
			flag.String("other", "", "other flag")

			// Set up args
			os.Args = tt.args

			// Parse flags
			flag.Parse()

			// Test
			got := IsFlagPassed(tt.flagName)
			if got != tt.want {
				t.Errorf("IsFlagPassed(%q) = %v, want %v", tt.flagName, got, tt.want)
			}
		})
	}
}
