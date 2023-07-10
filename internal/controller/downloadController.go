package controller

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/vpoluyaktov/abb_ia/internal/config"
	"github.com/vpoluyaktov/abb_ia/internal/dto"
	"github.com/vpoluyaktov/abb_ia/internal/ia_client"
	"github.com/vpoluyaktov/abb_ia/internal/logger"
	"github.com/vpoluyaktov/abb_ia/internal/mq"
	"github.com/vpoluyaktov/abb_ia/internal/utils"
)

type DownloadController struct {
	mq        *mq.Dispatcher
	item      *dto.IAItem
	startTime time.Time
	files     []fileDownload
	stopFlag  bool
}

type fileDownload struct {
	fileId          int
	fileSize        int64
	bytesDownloaded int64
	progress        int
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
	c.files = make([]fileDownload, len(c.item.Files))
	jd := utils.NewJobDispatcher(config.GetParallelDownloads())
	for i, f := range c.item.Files {
		jd.AddJob(i, ia.DownloadFile, outputDir, c.item.Server, c.item.Dir, f.Name, i, f.Size, c.updateFileProgress)
	}
	go c.updateTotalProgress()
	// if c.stopFlag {
	// 	break
	// }

	jd.Start()

	c.mq.SendMessage(mq.DownloadController, mq.Footer, &dto.SetBusyIndicator{Busy: false}, false)
	c.mq.SendMessage(mq.DownloadController, mq.Footer, &dto.UpdateStatus{Message: ""}, false)
	c.mq.SendMessage(mq.DownloadController, mq.DownloadPage, &dto.DownloadComplete{Audiobook: cmd.Audiobook}, true)
}

func (c *DownloadController) updateFileProgress(fileId int, fileName string, size int64, pos int64, percent int) {
	if c.files[fileId].progress != percent {
		// sent a message only if progress changed
		c.mq.SendMessage(mq.DownloadController, mq.DownloadPage, &dto.DownloadFileProgress{FileId: fileId, FileName: fileName, Percent: percent}, false)
	}
	c.files[fileId].fileId = fileId
	c.files[fileId].fileSize = size
	c.files[fileId].bytesDownloaded = pos
	c.files[fileId].progress = percent
}

func (c *DownloadController) updateTotalProgress() {
	var percent int = -1

	for !c.stopFlag && percent <= 100 {
		var totalSize = c.item.TotalSize
		var totalBytesDownloaded int64 = 0
		filesDownloaded := 0
		for _, f := range c.files {
			totalBytesDownloaded += f.bytesDownloaded
			if f.progress == 100 {
				filesDownloaded++
			}
		}

		var p int = 0
		if totalSize > 0 {
			p = int(float64(totalBytesDownloaded) / float64(totalSize) * 100)
		}

		// fix wrong file size returned by IA metadata
		if filesDownloaded == len(c.files) {
			p = 100
			totalBytesDownloaded = c.item.TotalSize
		}

		if percent != p {
			// sent a message only if progress changed
			percent = p

			elapsed := time.Since(c.startTime).Seconds()
			speed := int64(float64(totalBytesDownloaded) / elapsed)
			eta := (100 / (float64(percent) / elapsed)) - elapsed
			if eta < 0 || eta > (60*60*24*365) {
				eta = 0
			}

			elapsedH, _ := utils.SecondsToTime(elapsed)
			bytesH, _ := utils.BytesToHuman(totalBytesDownloaded)
			filesH := fmt.Sprintf("%d/%d", filesDownloaded, len(c.item.Files))
			speedH, _ := utils.SpeedToHuman(speed)
			etaH, _ := utils.SecondsToTime(eta)

			c.mq.SendMessage(mq.DownloadController, mq.DownloadPage, &dto.DownloadProgress{Elapsed: elapsedH, Percent: percent, Files: filesH, Bytes: bytesH, Speed: speedH, ETA: etaH}, false)
		}
		time.Sleep(mq.PullFrequency)
	}
}
