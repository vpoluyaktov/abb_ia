package ffmpeg

import (
	"os/exec"
	"strings"
)

type args struct {
	args []string
}

func NewArgs() *args {
	a := &args{}
	return a
}

func (a *args) Append(arg ...string) *args {
	for _, ar := range arg {
		ars := strings.Split(ar, " ")
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