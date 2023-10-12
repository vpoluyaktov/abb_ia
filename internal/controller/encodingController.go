package controller

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/vpoluyaktov/abb_ia/internal/config"
	"github.com/vpoluyaktov/abb_ia/internal/dto"
	"github.com/vpoluyaktov/abb_ia/internal/ffmpeg"
	"github.com/vpoluyaktov/abb_ia/internal/logger"
	"github.com/vpoluyaktov/abb_ia/internal/mq"
	"github.com/vpoluyaktov/abb_ia/internal/utils"
)

type EncodingController struct {
	mq        *mq.Dispatcher
	item      *dto.IAItem
	startTime time.Time
	files     []fileEncode
	stopFlag  bool
}

type fileEncode struct {
	fileId   int
	progress int
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
	logger.Info(mq.EncodingController + " received " + cmd.String())
	c.mq.SendMessage(mq.EncodingController, mq.Footer, &dto.UpdateStatus{Message: "Re-encoding mp3 files..."}, false)
	c.mq.SendMessage(mq.EncodingController, mq.Footer, &dto.SetBusyIndicator{Busy: true}, false)
	c.item = cmd.Audiobook.IAItem
	outputDir := filepath.Join("output", c.item.ID, c.item.Dir)

	// re-encode files
	c.stopFlag = false
	c.files = make([]fileEncode, len(c.item.Files))
	jd := utils.NewJobDispatcher(config.ParallelEncoders())
	for i, f := range c.item.Files {
		jd.AddJob(i, c.encodeFile, i, outputDir, f.Name)
	}
	go c.updateTotalProgress()
	// if c.stopFlag {
	// 	break
	// }

	jd.Start()

	c.mq.SendMessage(mq.EncodingController, mq.Footer, &dto.SetBusyIndicator{Busy: false}, false)
	c.mq.SendMessage(mq.EncodingController, mq.Footer, &dto.UpdateStatus{Message: ""}, false)
	c.mq.SendMessage(mq.EncodingController, mq.EncodingPage, &dto.EncodingComplete{Audiobook: cmd.Audiobook}, true)
}

func (c *EncodingController) encodeFile(fileId int, outputDir string, fileName string) {

	filePath := filepath.Join(outputDir, fileName)
	tmpFile := filePath + ".tmp"
	mp3, _ := ffmpeg.NewFFProbe(filePath)
	duration := mp3.Duration()

	// launch progress listener
	l, port := c.startProgressListener(fileId)
	defer l.Close()
	go c.updateFileProgress(fileId, fileName, duration, l)

	// launch ffmpeg process
	_, err := ffmpeg.NewFFmpeg().
		Input(filePath, "-f mp3").
		Output(tmpFile, fmt.Sprintf("-f mp3 -ab %s -ar %s -vn", config.BitRate(), config.SampleRate())).
		Overwrite(true).
		Params("-hide_banner -nostdin -nostats -loglevel fatal").
		SendProgressTo("http://127.0.0.1:" + strconv.Itoa(port)).
		Run()
	if err != nil {
		logger.Error("FFMPEG Error: " + string(err.Error()))
	} else {
		err := os.Remove(filePath)
		if err != nil {
			logger.Error("Can't delete file " + filePath + ": " + err.Error())
		} else {
			os.Rename(tmpFile, filePath)
		}
	}
}

func (c *EncodingController) startProgressListener(fileId int) (net.Listener, int) {

	basePortNumber := 31000
	portNumber := basePortNumber + fileId

	l, err := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(portNumber))
	if err != nil {
		logger.Error("Encoding progress: Start listener error: " + err.Error())
	}
	return l, portNumber
}

func (c *EncodingController) updateFileProgress(fileId int, fileName string, totalDuration float64, l net.Listener) {

	re := regexp.MustCompile(`out_time_ms=(\d+)`)
	fd, err := l.Accept()
	if err != nil {
		return // listener is closed
	}
	buf := make([]byte, 16)
	data := ""
	percent := 0
	for {
		_, err := fd.Read(buf)
		if err != nil {
			return // listener is closed
		}
		data += string(buf)
		a := re.FindAllStringSubmatch(data, -1)
		p := 0
		pstr := ""
		if len(a) > 0 && len(a[len(a)-1]) > 0 {
			c, _ := strconv.Atoi(a[len(a)-1][len(a[len(a)-1])-1])
			pstr = fmt.Sprintf("%.2f", float64(c)/totalDuration/1000000)
		}
		if strings.Contains(data, "progress=end") {
			p = 100
		}
		if pstr == "" {
			p = 0
		}
		pflt, err := strconv.ParseFloat(pstr, 64)
		if err != nil {
			p = 0
		} else {
			p = int(pflt * 100)
		}

		if p != percent {
			percent = p
			// sent a message only if progress changed
			c.mq.SendMessage(mq.EncodingController, mq.EncodingPage, &dto.EncodingFileProgress{FileId: fileId, FileName: fileName, Percent: percent}, true)
		}
		c.files[fileId].fileId = fileId
		c.files[fileId].progress = percent
	}
}

func (c *EncodingController) updateTotalProgress() {
	var percent int = -1

	for !c.stopFlag && percent <= 100 {
		var totalPercent int = 0
		files := 0
		for _, f := range c.files {
			totalPercent += f.progress
			if f.progress == 100 {
				files++
			}
		}
		p := int(totalPercent / len(c.files))

		if percent != p {
			// sent a message only if progress changed
			percent = p

			elapsed := time.Since(c.startTime).Seconds()
			speed := int64(float64(percent) / elapsed)
			eta := (100 / (float64(percent) / elapsed)) - elapsed
			if eta < 0 || eta > (60*60*24*365) {
				eta = 0
			}

			elapsedH, _ := utils.SecondsToTime(elapsed)
			filesH := fmt.Sprintf("%d/%d", files, len(c.item.Files))
			speedH, _ := utils.SpeedToHuman(speed)
			etaH, _ := utils.SecondsToTime(eta)

			c.mq.SendMessage(mq.EncodingController, mq.EncodingPage, &dto.EncodingProgress{Elapsed: elapsedH, Percent: percent, Files: filesH, Speed: speedH, ETA: etaH}, true)
		}
		time.Sleep(mq.PullFrequency)
	}
}
