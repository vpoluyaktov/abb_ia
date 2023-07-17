package ffmpeg

import (
	"os/exec"

	"github.com/vpoluyaktov/abb_ia/internal/logger"
)

type FFmpeg struct {
	input  input
	output output
	params params
}

type input struct {
	fileName string
	args     string
}

type output struct {
	fileName string
	args     string
}

type params struct {
	args string
}

func NewFFmpeg() *FFmpeg {
	ffmpeg := &FFmpeg{
		input:  input{},
		output: output{},
		params: params{},
	}
	return ffmpeg
}

func (f *FFmpeg) Input(fileName string, args string) *FFmpeg {
	f.input.fileName = fileName
	f.input.args = args
	return f
}

func (f *FFmpeg) Output(fileName string, args string) *FFmpeg {
	f.output.fileName = fileName
	f.output.args = args
	return f
}

func (f *FFmpeg) Params(args string) *FFmpeg {
	f.params.args += " " + args
	return f
}

func (f *FFmpeg) SendProgressTo(url string) *FFmpeg {
	f.params.args += " -progress " + url
	return f
}

func (f *FFmpeg) Overwrite(b bool) *FFmpeg {
	if b {
		f.params.args += " -y"
	}
	return f
}

func (f *FFmpeg) Run() (string, *exitErr) {
	cmd := "ffmpeg"
	args := NewArgs().
		AppendArgs(f.params.args).
		AppendArgs(f.input.args).AppendArgs("-i").AppendFileName(f.input.fileName).
		AppendArgs(f.output.args).AppendFileName(f.output.fileName)
	command := exec.Command(cmd, args.String()...)
	logger.Debug("FFMPEG cmd: " + command.String())
	out, err := command.Output()
	if err != nil {
		return string(out), ExitErr(err)
	} else {
		return string(out), nil
	}
}
