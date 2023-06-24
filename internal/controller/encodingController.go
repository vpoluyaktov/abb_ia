package controller

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	ffmpeg "github.com/u2takey/ffmpeg-go"

	"github.com/vpoluyaktov/audiobook_creator_IA/internal/dto"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/logger"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/mq"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/utils"
)

type EncodingController struct {
	mq         *mq.Dispatcher
	item       *dto.IAItem
	startTime  time.Time
	progress   []int
	downloaded []int64
	stopFlag   bool
}

func NewEncodingController(dispatcher *mq.Dispatcher) *EncodingController {
	dc := &EncodingController{}
	dc.mq = dispatcher
	dc.mq.RegisterListener(mq.EncodingController, dc.dispatchMessage)
	return dc
}

func (c *EncodingController) checkMQ() {
	m := c.mq.GetMessage(mq.EncodingController)
	if m != nil {
		c.dispatchMessage(m)
	}
}

func (c *EncodingController) dispatchMessage(m *mq.Message) {
	switch dto := m.Dto.(type) {
	case *dto.EncodeCommand:
		go c.startEncoding(dto)
	case *dto.StopCommand:
		go c.stopEncoding(dto)
	default:
		m.UnsupportedTypeError(mq.EncodingController)
	}
}

func (c *EncodingController) stopEncoding(cmd *dto.StopCommand) {
	c.stopFlag = true
	logger.Debug(mq.EncodingController + ": Received StopEncoding command")
}

func (c *EncodingController) startEncoding(cmd *dto.EncodeCommand) {
	c.startTime = time.Now()
	logger.Debug(mq.EncodingController + ": Received StartEncoding command with IA item: " + cmd.String())
	c.mq.SendMessage(mq.EncodingController, mq.Footer, &dto.UpdateStatus{Message: "Re-encoding mp3 files..."}, false)
	c.mq.SendMessage(mq.EncodingController, mq.Footer, &dto.SetBusyIndicator{Busy: true}, false)
	c.item = cmd.Audiobook.IAItem
	outputDir := filepath.Join("output", c.item.ID, c.item.Dir)

	// re-encode files
	c.stopFlag = false
	c.progress = make([]int, len(c.item.Files))
	c.downloaded = make([]int64, len(c.item.Files))
	go c.updateEncodingProgress()
	for i, f := range c.item.Files {
		if c.stopFlag {
			break
		}
		c.encodeFile(i, outputDir, f.Name, c.updateFileProgress)
	}
	c.mq.SendMessage(mq.EncodingController, mq.Footer, &dto.SetBusyIndicator{Busy: false}, false)
	c.mq.SendMessage(mq.EncodingController, mq.Footer, &dto.UpdateStatus{Message: ""}, false)
}

type Fn func(int, string, int64, int)

func (c *EncodingController) encodeFile(id int, dir string, file string, updateProgress Fn) {

	filePath := filepath.Join(dir, file)
	tmpFile := filePath + ".tmp"
	a, err := ffmpeg.Probe(filePath)
	if err != nil {
		logger.Error("FFMPEG Probe Error: " + err.Error())
	}
	totalDuration, err := probeDuration(a)
	if err != nil {
		logger.Error("FFMPEG ProbeDuration Error: " + err.Error())
	}
	TempSock(totalDuration)
	err = ffmpeg.Input(filePath).
		Output(tmpFile, ffmpeg.KwArgs{"c:v": "libx264", "preset": "veryslow", "f": "mp3"}).
		GlobalArgs("-progress", "http://127.0.0.1:31001").
		OverWriteOutput().
		Run()
	if err != nil {
		logger.Error("FFMPEG Error: " + err.Error())
	}
}

func (c *EncodingController) updateFileProgress(fileId int, fileName string, pos int64, percent int) {
	if c.progress[fileId] != percent {
		// sent a message only if progress changed
		c.mq.SendMessage(mq.EncodingController, mq.EncodingPage, &dto.EncodingFileProgress{FileId: fileId, FileName: fileName, Percent: percent}, true)
	}
	c.progress[fileId] = percent
	c.downloaded[fileId] = pos
}

func (c *EncodingController) updateEncodingProgress() {
	var percent int = -1
	var files int = 0
	var speed int64 = 0
	var eta float64 = 0
	var bytes int64 = 0

	for !c.stopFlag && percent <= 100 {
		var totalPercent int = 0
		files = 0
		for _, p := range c.progress {
			totalPercent += p
			if p == 100 {
				files++
			}
		}
		p := int(totalPercent / len(c.progress))

		if percent != p {
			// sent a message only if progress changed
			percent = p

			bytes = 0
			for _, b := range c.downloaded {
				bytes += b
			}

			duration := time.Since(c.startTime).Seconds()
			speed = int64(float64(bytes) / duration)
			eta = (100 / (float64(percent) / duration)) - duration
			if eta < 0 || eta > (60*60*24*365) {
				eta = 0
			}

			durationH, _ := utils.SecondsToTime(duration)
			bytesH, _ := utils.BytesToHuman(bytes)
			filesH := fmt.Sprintf("%d/%d", files, len(c.item.Files))
			speedH, _ := utils.SpeedToHuman(speed)
			etaH, _ := utils.SecondsToTime(eta)

			c.mq.SendMessage(mq.EncodingController, mq.EncodingPage, &dto.EncodingProgress{Duration: durationH, Percent: percent, Files: filesH, Bytes: bytesH, Speed: speedH, ETA: etaH}, true)
		}
		time.Sleep(mq.PullFrequency)
	}

	// c.mq.SendMessage(mq.EncodingController, mq.EncodingPage, &dto.EncodingComplete{}, true)
}

func TempSock(totalDuration float64) string {
	// serve

	rand.Seed(time.Now().Unix())
	sockFileName := path.Join(os.TempDir(), fmt.Sprintf("%d_sock", rand.Int()))
	l, err := net.Listen("tcp", "127.0.0.1:31001")
	if err != nil {
		logger.Error("Encoding progress listener Error: " + err.Error())
	}

	go func() {
		re := regexp.MustCompile(`out_time_ms=(\d+)`)
		fd, err := l.Accept()
		if err != nil {
			logger.Error("Encoding progress listener Error: " + err.Error())
		}
		buf := make([]byte, 16)
		data := ""
		progress := ""
		for {
			_, err := fd.Read(buf)
			if err != nil {
				return
			}
			data += string(buf)
			a := re.FindAllStringSubmatch(data, -1)
			cp := ""
			if len(a) > 0 && len(a[len(a)-1]) > 0 {
				c, _ := strconv.Atoi(a[len(a)-1][len(a[len(a)-1])-1])
				cp = fmt.Sprintf("%.2f", float64(c)/totalDuration/1000000)
			}
			if strings.Contains(data, "progress=end") {
				cp = "done"
			}
			if cp == "" {
				cp = ".0"
			}
			if cp != progress {
				progress = cp
				logger.Debug("Encoding progress: " + progress)
			}
		}
	}()

	return sockFileName
}

type probeFormat struct {
	Duration string `json:"duration"`
}

type probeData struct {
	Format probeFormat `json:"format"`
}

func probeDuration(a string) (float64, error) {
	pd := probeData{}
	err := json.Unmarshal([]byte(a), &pd)
	if err != nil {
		return 0, err
	}
	f, err := strconv.ParseFloat(pd.Format.Duration, 64)
	if err != nil {
		return 0, err
	}
	return f, nil
}
