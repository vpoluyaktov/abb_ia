package ffmpeg

import (
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type args struct {
	args []string
}

func NewArgs() *args {
	a := &args{}
	return a
}

func (a *args) AppendFileName(arg string) *args {
	a.args = append(a.args, arg)
	return a
}

func (a *args) AppendArgs(arg ...string) *args {
	for _, ar := range arg {
		ars := strings.Fields(ar)
		a.args = append(a.args, ars...)
	}
	return a
}

func (a *args) String() []string {
	return a.args
}

type exitErr struct {
	err        error
	isExitErr  bool
	errMessage string
}

func ExitErr(e error) *exitErr {
	ee := &exitErr{
		err: e,
	}
	exitErr, isExitError := e.(*exec.ExitError)
	if isExitError {
		ee.isExitErr = true
		ee.errMessage = string(exitErr.Stderr)
	} else {
		ee.errMessage = e.Error()
	}
	return ee
}

func (e *exitErr) Error() string {
	return e.errMessage
}

// parse ffmpeg progress stats
var reTotalSize = regexp.MustCompile(`total_size=(\d+)`)
var reOutTimeUs = regexp.MustCompile(`out_time_us=(\d+)`)
var reSpeed = regexp.MustCompile(`speed=\s*(\d+)x`)
var reComplete = regexp.MustCompile(`progress=end`)

func ParseFFMPEGProgress(data string) (int64, float64, float64, bool) {

	var bytesProcessed int64 = 0
	var secondsProcessed float64 = 0
	var encodingSpeed float64 = 0
	var complete = false

	totalSizeMatches := reTotalSize.FindAllStringSubmatch(data, -1)
	if len(totalSizeMatches) > 0 {
		lastTotalSize := totalSizeMatches[len(totalSizeMatches)-1][1]
		bytesProcessed, _ = strconv.ParseInt(lastTotalSize, 10, 64)
	}

	outTimeUsMatches := reOutTimeUs.FindAllStringSubmatch(data, -1)
	if len(outTimeUsMatches) > 0 {
		lastOutTimeUs := outTimeUsMatches[len(outTimeUsMatches)-1][1]
		msecProcessed, _ := strconv.ParseFloat(lastOutTimeUs, 64)
		secondsProcessed = msecProcessed / 1000000
	}

	speedMatches := reSpeed.FindAllStringSubmatch(data, -1)
	if len(speedMatches) > 0 {
		lastSpeed := speedMatches[len(speedMatches)-1][1]
		encodingSpeed, _ = strconv.ParseFloat(lastSpeed, 64)

	}

	completeMatches := reComplete.FindAllStringSubmatch(data, -1)
	if len(completeMatches) > 0 {
		complete = true
	}

	return bytesProcessed, secondsProcessed, encodingSpeed, complete
}
