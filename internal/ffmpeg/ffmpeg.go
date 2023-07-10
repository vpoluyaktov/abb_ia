package ffmpeg

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
	f.params.args = args
	return f
}

func (f *FFmpeg) Run() error {

	var err error
	return err
}
