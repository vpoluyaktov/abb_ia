package utils_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/logger"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/utils"
)

const (
	logFileName string = "/tmp/audiobook_creator_IA.test.log"
	logLevel           = logger.DEBUG
)

func TestMain(m *testing.M) {
	logger.Init(logFileName, logLevel)
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