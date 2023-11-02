package utils_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vpoluyaktov/abb_ia/internal/config"
	"github.com/vpoluyaktov/abb_ia/internal/logger"
	"github.com/vpoluyaktov/abb_ia/internal/utils"
)

func TestMain(m *testing.M) {
	config.Load()
	logger.Init(config.LogFileName(), "DEBUG")
	os.Exit(m.Run())
}

func TestTimeToSeconds(t *testing.T) {
	tStr := "1:04:05"
	tInt, err := utils.TimeToSeconds(tStr)
	assert.NoError(t, err)
	assert.Equal(t, float64(3845), tInt)

	tStr = "1:4:5"
	tInt, err = utils.TimeToSeconds(tStr)
	assert.NoError(t, err)
	assert.Equal(t, float64(3845), tInt)

	tStr = "04:05"
	tInt, err = utils.TimeToSeconds(tStr)
	assert.NoError(t, err)
	assert.Equal(t, float64(245), tInt)

	tStr = "05"
	tInt, err = utils.TimeToSeconds(tStr)
	assert.NoError(t, err)
	assert.Equal(t, float64(5), tInt)

	tStr = "3845.32"
	tInt, err = utils.TimeToSeconds(tStr)
	assert.NoError(t, err)
	assert.Equal(t, float64(3845.32), tInt)

}

func TestSecondToTime(t *testing.T) {
	sec := float64((5 * 3600) + (45 * 60) + 35)
	time := utils.SecondsToTime(sec)
	assert.Equal(t, " 5:45:35", time)
}

func TestBytesToHuman(t *testing.T) {
	b := int64((5 * 1024 * 1024) + (245 * 1024) + 35)
	size := utils.BytesToHuman(b)
	assert.Equal(t, "5.2 Mb", size)
}

func TestHumanToBytes(t *testing.T) {
	// Test cases with valid inputs
	validInputs := []struct {
			input       string
			expected    int64
	}{
			{"100 b", 100},
			{"150B", 150},
			{"1 kb", 1024},
			{"1 Kb", 1024},
			{"2 mb", 2097152},
			{"2Mb", 2097152},
			{"155.5 Mb", 163053568},
			{"155.5Mb", 163053568},		
			{"3 gb", 3221225472},
			{"3 Gb", 3221225472},
	}

	for _, tc := range validInputs {
			actual, err := utils.HumanToBytes(tc.input)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, actual)
	}
}