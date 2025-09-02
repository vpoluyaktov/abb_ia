package controller

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"abb_ia/internal/config"
	"abb_ia/internal/dto"
	"abb_ia/internal/ffmpeg"
	"abb_ia/internal/logger"
	"abb_ia/internal/mq"
	"abb_ia/internal/utils"
)

type EncodingController struct {
	mq        *mq.Dispatcher
	ab        *dto.Audiobook
	startTime time.Time
	files     []fileEncode
	stopFlag  bool
}

// progress tracking arrays
type fileEncode struct {
	fileId           int
	fileName         string
	filePath         string
	totalDuration    float64
	bytesProcessed   int64
	secondsProcessed float64
	encodingSpeed    float64
	progress         int
	complete         bool
}

func NewEncodingController(dispatcher *mq.Dispatcher) *EncodingController {
	c := &EncodingController{}
	c.mq = dispatcher
	c.mq.RegisterListener(mq.EncodingController, c.dispatchMessage)
	return c
}

func (c *EncodingController) checkMQ() {
	m, err := c.mq.GetMessage(mq.EncodingController)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to get message for EncodingController: %v", err))
		return
	}
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

func (c *EncodingController) stopEncoding(_ *dto.StopCommand) {
	c.stopFlag = true
	logger.Debug(mq.EncodingController + ": Received StopEncoding command")
}

func (c *EncodingController) startEncoding(cmd *dto.EncodeCommand) {
	c.startTime = time.Now()

	c.ab = cmd.Audiobook
	c.stopFlag = false
	c.files = make([]fileEncode, len(c.ab.Mp3Files))

	c.mq.SendMessage(mq.EncodingController, mq.EncodingPage, &dto.DisplayBookInfoCommand{Audiobook: c.ab}, mq.PriorityHigh)
	c.mq.SendMessage(mq.EncodingController, mq.Footer, &dto.UpdateStatus{Message: "Re-encoding mp3 files..."}, mq.PriorityNormal)
	c.mq.SendMessage(mq.EncodingController, mq.Footer, &dto.SetBusyIndicator{Busy: true}, mq.PriorityNormal)

	logger.Info(fmt.Sprintf("Re-encoding mp3 files: %s - %s...", c.ab.Author, c.ab.Title))

	// re-encode files
	jd := utils.NewJobDispatcher(c.ab.Config.GetConcurrentEncoders())
	for i, f := range c.ab.Mp3Files {
		c.files[i].fileId = i
		c.files[i].fileName = f.FileName
		c.files[i].filePath = filepath.Join(c.ab.OutputDir, f.FileName)
		mp3, _ := ffmpeg.NewFFProbe(c.files[i].filePath)
		c.files[i].totalDuration = mp3.Duration()

		jd.AddJob(i, c.encodeFile, i, c.ab.OutputDir)
	}
	go c.updateTotalProgress()
	jd.Start()

	c.mq.SendMessage(mq.EncodingController, mq.Footer, &dto.SetBusyIndicator{Busy: false}, mq.PriorityNormal)
	c.mq.SendMessage(mq.EncodingController, mq.Footer, &dto.UpdateStatus{Message: ""}, mq.PriorityNormal)
	if !c.stopFlag {
		c.mq.SendMessage(mq.EncodingController, mq.EncodingPage, &dto.EncodingComplete{Audiobook: cmd.Audiobook}, mq.PriorityHigh)
	}
	c.stopFlag = true
}

func (c *EncodingController) encodeFile(fileId int, outputDir string) {
	if c.stopFlag {
		return
	}

	filePath := c.files[fileId].filePath
	tmpFile := filePath + ".tmp"

	// launch progress listener
	l, port := c.startProgressListener(fileId)
	defer l.Close()
	go c.updateFileProgress(fileId, l)

	// launch ffmpeg process
	ffmpeg := ffmpeg.NewFFmpeg().
		Input(filePath, "-f mp3").
		Output(tmpFile, fmt.Sprintf("-f mp3 -ab %dk -ar %d -vn", c.ab.Config.GetBitRate(), c.ab.Config.GetSampleRate())).
		Overwrite(true).
		Params("-hide_banner -nostdin -nostats -loglevel error").
		SendProgressTo("http://127.0.0.1:" + strconv.Itoa(port))

	go c.killSwitch(ffmpeg)
	_, err := ffmpeg.Run()
	if err != nil && !c.stopFlag {
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

func (c *EncodingController) killSwitch(ffmpeg *ffmpeg.FFmpeg) {
	for !c.stopFlag {
		time.Sleep(mq.PullFrequency)
	}
	ffmpeg.Kill()
}

func (c *EncodingController) startProgressListener(fileId int) (net.Listener, int) {
	portNumber := config.Instance().GetBasePortNumber() + fileId
	l, err := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(portNumber))
	if err != nil {
		logger.Error("Encoding progress: Start listener error: " + err.Error())
	}
	return l, portNumber
}

func (c *EncodingController) updateFileProgress(fileId int, l net.Listener) {
	fd, err := l.Accept()
	if err != nil {
		return // listener may be closed already
	}
	buf := make([]byte, 16)
	data := ""
	percent := 0

	for !c.stopFlag {
		_, err := fd.Read(buf)
		if err != nil {
			return // listener is closed
		}
		data += string(buf)
		bytesProcessed, secondsProcessed, encodingSpeed, complete := ffmpeg.ParseFFMPEGProgress(data)
		percent = int(secondsProcessed / c.files[fileId].totalDuration * 100)
		// wrong calculation protection
		if percent < 0 {
			percent = 0
		} else if percent > 100 {
			percent = 100
		} else if percent < c.files[fileId].progress {
			percent = c.files[fileId].progress
		} else if complete {
			percent = 100
		}

		// sent a message only if progress changed
		if percent != c.files[fileId].progress {
			c.files[fileId].bytesProcessed = bytesProcessed
			c.files[fileId].secondsProcessed = secondsProcessed
			c.files[fileId].encodingSpeed = encodingSpeed
			c.files[fileId].progress = percent
			c.files[fileId].complete = complete
			c.mq.SendMessage(mq.EncodingController, mq.EncodingPage, &dto.EncodingFileProgress{FileId: fileId, FileName: c.files[fileId].fileName, Percent: percent}, mq.PriorityHigh)
		}
	}
}

func (c *EncodingController) updateTotalProgress() {
	var p int = -1

	for !c.stopFlag {
		var totalDuration float64 = 0
		var secondsProcessed float64 = 0
		var totalSpeed float64 = 0
		filesProcessed := 0
		filesComplete := 0
		for _, f := range c.files {
			totalDuration += f.totalDuration
			secondsProcessed += f.secondsProcessed
			totalSpeed += f.encodingSpeed
			if f.complete {
				filesComplete++
			}
			if f.encodingSpeed > 0 {
				filesProcessed++
			}
		}
		percent := int(secondsProcessed / totalDuration * 100)
		// wrong calculation protection
		if percent < 0 {
			percent = 0
		} else if percent > 100 {
			percent = 100
		}

		if percent != p {
			// sent a message only if progress changed
			p = percent

			elapsed := time.Since(c.startTime).Seconds()
			speed := totalSpeed / float64(filesProcessed)
			eta := (100 / (float64(percent) / elapsed)) - elapsed
			if eta < 0 || eta > (60*60*24*365) {
				eta = 0
			}

			elapsedH := utils.SecondsToTime(elapsed)
			filesH := fmt.Sprintf("%d/%d", filesComplete, len(c.ab.Mp3Files))
			speedH := fmt.Sprintf("%.0fx", speed)
			etaH := utils.SecondsToTime(eta)

			c.mq.SendMessage(mq.EncodingController, mq.EncodingPage, &dto.EncodingProgress{Elapsed: elapsedH, Percent: percent, Files: filesH, Speed: speedH, ETA: etaH}, mq.PriorityHigh)
		}
		time.Sleep(mq.PullFrequency)
	}
}
