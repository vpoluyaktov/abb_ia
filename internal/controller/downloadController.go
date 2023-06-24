package controller

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/vpoluyaktov/audiobook_creator_IA/internal/config"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/dto"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/ia_client"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/logger"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/mq"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/utils"
)

type DownloadController struct {
	mq         *mq.Dispatcher
	item       *dto.IAItem
	startTime  time.Time
	progress   []int
	downloaded []int64
	stopFlag   bool
}

func NewDownloadController(dispatcher *mq.Dispatcher) *DownloadController {
	dc := &DownloadController{}
	dc.mq = dispatcher
	dc.mq.RegisterListener(mq.DownloadController, dc.dispatchMessage)
	return dc
}

func (c *DownloadController) checkMQ() {
	m := c.mq.GetMessage(mq.DownloadController)
	if m != nil {
		c.dispatchMessage(m)
	}
}

func (c *DownloadController) dispatchMessage(m *mq.Message) {
	switch dto := m.Dto.(type) {
	case *dto.DownloadCommand:
		go c.startDownload(dto)
	case *dto.StopCommand:
		go c.stopDownload(dto)
	default:
		m.UnsupportedTypeError(mq.DownloadController)
	}
}

func (c *DownloadController) stopDownload(cmd *dto.StopCommand) {
	c.stopFlag = true
	logger.Debug(mq.DownloadController + ": Received StopDownload command")
}

func (c *DownloadController) startDownload(cmd *dto.DownloadCommand) {
	c.startTime = time.Now()
	logger.Debug(mq.DownloadController + ": Received StartDownload command with IA item: " + cmd.String())
	c.mq.SendMessage(mq.DownloadController, mq.Footer, &dto.UpdateStatus{Message: "Downloading mp3 files..."}, false)
	c.mq.SendMessage(mq.DownloadController, mq.Footer, &dto.SetBusyIndicator{Busy: true}, false)
	c.item = cmd.Audiobook.IAItem
	outputDir := filepath.Join("output", c.item.ID)

	// download files
	ia := ia_client.New(config.IsUseMock(), config.IsSaveMock())
	c.stopFlag = false
	c.progress = make([]int, len(c.item.Files))
	c.downloaded = make([]int64, len(c.item.Files))
	jd := utils.NewJobDispatcher(config.GetParallelDownloads())
	for i, f := range c.item.Files {
		jd.AddJob(i, ia.DownloadFile, outputDir, c.item.Server, c.item.Dir, f.Name, i, c.updateFileProgress)
	}
	// if c.stopFlag {
	// 	break
	// }
	go c.updateDownloadProgress()
	jd.Start()
	c.mq.SendMessage(mq.DownloadController, mq.Footer, &dto.SetBusyIndicator{Busy: false}, false)
	c.mq.SendMessage(mq.DownloadController, mq.Footer, &dto.UpdateStatus{Message: ""}, false)
	c.mq.SendMessage(mq.DownloadController, mq.DownloadPage, &dto.DownloadComplete{Audiobook: cmd.Audiobook}, true)
}

func (c *DownloadController) updateFileProgress(fileId int, fileName string, pos int64, percent int) {
	if c.progress[fileId] != percent {
		// sent a message only if progress changed
		c.mq.SendMessage(mq.DownloadController, mq.DownloadPage, &dto.DownloadFileProgress{FileId: fileId, FileName: fileName, Percent: percent}, true)
	}
	c.progress[fileId] = percent
	c.downloaded[fileId] = pos
}

func (c *DownloadController) updateDownloadProgress() {
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
			if eta < 0 || eta > (60 * 60 * 24 * 365) {
				eta = 0
			}

			durationH, _ := utils.SecondsToTime(duration)
			bytesH, _ := utils.BytesToHuman(bytes)
			filesH := fmt.Sprintf("%d/%d", files, len(c.item.Files))
			speedH, _ := utils.SpeedToHuman(speed)
			etaH, _ := utils.SecondsToTime(eta)

			c.mq.SendMessage(mq.DownloadController, mq.DownloadPage, &dto.DownloadProgress{Duration: durationH, Percent: percent, Files: filesH, Bytes: bytesH, Speed: speedH, ETA: etaH}, false)
		}
		time.Sleep(mq.PullFrequency)
	}
}
