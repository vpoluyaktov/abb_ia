package utils

import (
	"testing"
)

func TestTimeToSeconds(t *testing.T) {
	tests := []struct {
		name    string
		time    string
		want    float64
		wantErr bool
	}{
		{
			name:    "empty string",
			time:    "",
			want:    0,
			wantErr: true,
		},
		{
			name:    "invalid format",
			time:    "abc",
			want:    0,
			wantErr: true,
		},
		{
			name:    "HH:MM:SS format",
			time:    "01:30:45",
			want:    5445,
			wantErr: false,
		},
		{
			name:    "MM:SS format",
			time:    "30:45",
			want:    1845,
			wantErr: false,
		},
		{
			name:    "seconds only",
			time:    "45",
			want:    45,
			wantErr: false,
		},
		{
			name:    "seconds with milliseconds",
			time:    "45.5",
			want:    45.5,
			wantErr: false,
		},
		{
			name:    "invalid HH:MM:SS",
			time:    "01:aa:45",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TimeToSeconds(tt.time)
			if (err != nil) != tt.wantErr {
				t.Errorf("TimeToSeconds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TimeToSeconds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSecondsToTime(t *testing.T) {
	tests := []struct {
		name     string
		seconds  float64
		expected string
	}{
		{
			name:     "zero seconds",
			seconds:  0,
			expected: "0:00:00",
		},
		{
			name:     "seconds only",
			seconds:  45,
			expected: "0:00:45",
		},
		{
			name:     "minutes and seconds",
			seconds:  125,
			expected: "0:02:05",
		},
		{
			name:     "hours, minutes, and seconds",
			seconds:  3725,
			expected: "1:02:05",
		},
		{
			name:     "multiple hours",
			seconds:  7325,
			expected: "2:02:05",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SecondsToTime(tt.seconds)
			if got != tt.expected {
				t.Errorf("SecondsToTime() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestBytesToHuman(t *testing.T) {
	tests := []struct {
		name     string
		bytes    int64
		expected string
	}{
		{
			name:     "bytes",
			bytes:    500,
			expected: "500 b",
		},
		{
			name:     "kilobytes",
			bytes:    1500,
			expected: "1.5 Kb",
		},
		{
			name:     "megabytes",
			bytes:    1500000,
			expected: "1.4 Mb",
		},
		{
			name:     "gigabytes",
			bytes:    1500000000,
			expected: "1.4 Gb",
		},
		{
			name:     "terabytes",
			bytes:    1500000000000,
			expected: "1.4 Tb",
		},
		{
			name:     "petabytes",
			bytes:    1500000000000000,
			expected: "1.3 Pb",
		},
		{
			name:     "exabytes",
			bytes:    1500000000000000000,
			expected: "1.3 Eb",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BytesToHuman(tt.bytes)
			if got != tt.expected {
				t.Errorf("BytesToHuman() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestHumanToBytes(t *testing.T) {
	tests := []struct {
		name     string
		size     string
		expected int64
		wantErr  bool
	}{
		{
			name:     "bytes",
			size:     "500 b",
			expected: 500,
			wantErr:  false,
		},
		{
			name:     "kilobytes",
			size:     "1.5 kb",
			expected: 1536,
			wantErr:  false,
		},
		{
			name:     "megabytes",
			size:     "1.5 mb",
			expected: 1572864,
			wantErr:  false,
		},
		{
			name:     "gigabytes",
			size:     "1.5 gb",
			expected: 1610612736,
			wantErr:  false,
		},
		{
			name:     "terabytes",
			size:     "1.5 tb",
			expected: 1649267441664,
			wantErr:  false,
		},
		{
			name:     "concatenated format",
			size:     "1.5kb",
			expected: 1536,
			wantErr:  false,
		},
		{
			name:     "invalid number",
			size:     "abc kb",
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "invalid suffix",
			size:     "1.5 xb",
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "empty string",
			size:     "",
			expected: 0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HumanToBytes(tt.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("HumanToBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("HumanToBytes() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSeparateNumberAndSuffix(t *testing.T) {
	tests := []struct {
		name         string
		size         string
		wantNum     string
		wantSuffix string
	}{
		{
			name:         "simple format",
			size:         "1.5kb",
			wantNum:     "1.5",
			wantSuffix: "kb",
		},
		{
			name:         "integer format",
			size:         "500mb",
			wantNum:     "500",
			wantSuffix: "mb",
		},
		{
			name:         "with space",
			size:         "1.5 kb",
			wantNum:     "1.5",
			wantSuffix: "kb",
		},
		{
			name:         "empty string",
			size:         "",
			wantNum:     "",
			wantSuffix: "",
		},
		{
			name:         "only number",
			size:         "1.5",
			wantNum:     "1.5",
			wantSuffix: "",
		},
		{
			name:         "only suffix",
			size:         "kb",
			wantNum:     "",
			wantSuffix: "kb",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNum, gotSuffix := separateNumberAndSuffix(tt.size)
			if gotNum != tt.wantNum {
				t.Errorf("separateNumberAndSuffix() gotNum = %v, want %v", gotNum, tt.wantNum)
			}
			if gotSuffix != tt.wantSuffix {
				t.Errorf("separateNumberAndSuffix() gotSuffix = %v, want %v", gotSuffix, tt.wantSuffix)
			}
		})
	}
}

func TestSpeedToHuman(t *testing.T) {
	tests := []struct {
		name     string
		bytes    int64
		expected string
	}{
		{
			name:     "bytes per second",
			bytes:    500,
			expected: "500 b/s",
		},
		{
			name:     "kilobytes per second",
			bytes:    1500,
			expected: "1.5 Kb/s",
		},
		{
			name:     "megabytes per second",
			bytes:    1500000,
			expected: "1.4 Mb/s",
		},
		{
			name:     "gigabytes per second",
			bytes:    1500000000,
			expected: "1.4 Gb/s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SpeedToHuman(tt.bytes)
			if got != tt.expected {
				t.Errorf("SpeedToHuman() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFirstN(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		n        int
		expected string
	}{
		{
			name:     "empty string",
			str:      "",
			n:        5,
			expected: "",
		},
		{
			name:     "string shorter than n",
			str:      "hello",
			n:        10,
			expected: "hello",
		},
		{
			name:     "string longer than n",
			str:      "hello world",
			n:        5,
			expected: "hello…",
		},
		{
			name:     "n equals string length",
			str:      "hello",
			n:        5,
			expected: "hello",
		},
		{
			name:     "n is zero",
			str:      "hello",
			n:        0,
			expected: "…",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FirstN(tt.str, tt.n)
			if got != tt.expected {
				t.Errorf("FirstN() = %v, want %v", got, tt.expected)
			}
		})
	}
}
