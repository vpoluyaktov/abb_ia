package ffmpeg

import (
	"encoding/json"
	"os/exec"
	"path/filepath"
	"strconv"
)

type FFProbe struct {
	fileName string
	metadata Metadata
}

type Metadata struct {
	Format struct {
		Filename       string `json:"filename"`
		NbStreams      int    `json:"nb_streams"`
		NbPrograms     int    `json:"nb_programs"`
		FormatName     string `json:"format_name"`
		FormatLongName string `json:"format_long_name"`
		StartTime      string `json:"start_time"`
		Duration       string `json:"duration"`
		Size           string `json:"size"`
		BitRate        string `json:"bit_rate"`
		ProbeScore     int    `json:"probe_score"`
		Tags           struct {
			TLEN     string `json:"TLEN"`
			Title    string `json:"title"`
			Artist   string `json:"artist"`
			Album    string `json:"album"`
			TIT3     string `json:"TIT3"`
			Comment  string `json:"comment"`
			ITunNORM string `json:"iTunNORM"`
			Genre    string `json:"genre"`
			Date     string `json:"date"`
		} `json:"tags"`
	} `json:"format"`
}

func NewFFProbe(fileName string) (*FFProbe, error) {
	p := &FFProbe{fileName, Metadata{}}
	cmd := "ffprobe"
	args := NewArgs().
		AppendArgs("-loglevel error").
		AppendArgs("-show_format").
		AppendArgs("-show_streams").
		AppendArgs("-of json").
		AppendFileName(p.fileName)
	out, err := exec.Command(cmd, args.String()...).Output()
	if err == nil {
		err = json.Unmarshal([]byte(out), &p.metadata)
	}
	return p, err
}

func (p *FFProbe) GetDuration() float64 {
	f, err := strconv.ParseFloat(p.metadata.Format.Duration, 64)
	if err != nil {
		return 0
	}
	return f
}

func (p *FFProbe) GetTitle() string {
	if p.metadata.Format.Tags.Title == "" {
		return filepath.Base(p.metadata.Format.Filename)
	} else {
		return p.metadata.Format.Tags.Title
	}
}
