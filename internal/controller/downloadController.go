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
	ab        *dto.Audiobook
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
	logger.Info(mq.DownloadController + " received " + cmd.String())
	c.mq.SendMessage(mq.DownloadController, mq.Footer, &dto.UpdateStatus{Message: "Downloading mp3 files..."}, false)
	c.mq.SendMessage(mq.DownloadController, mq.Footer, &dto.SetBusyIndicator{Busy: true}, false)

	c.ab = cmd.Audiobook
	item := c.ab.IAItem
	c.ab.Author = item.Creator
	c.ab.Title = item.Title
	c.ab.Description = item.Description
	c.ab.CoverURL = item.Cover
	c.ab.OutputDir = utils.SanitizeFilePath(filepath.Join(config.OutputDir(), item.ID))
	c.ab.TotalSize = item.TotalSize
	c.ab.TotalDuration = item.TotalLength

	// update Book info on UI
	c.mq.SendMessage(mq.DownloadController, mq.DownloadPage, &dto.DisplayBookInfoCommand{Audiobook: c.ab}, true)

	// download files
	ia := ia_client.New(config.IsUseMock(), config.IsSaveMock())
	c.stopFlag = false
	c.files = make([]fileDownload, len(item.Files))
	jd := utils.NewJobDispatcher(config.ParallelDownloads())
	for i, iaFile := range item.Files {
		localFileName := utils.SanitizeFilePath(filepath.Join(item.Dir, iaFile.Name))
		c.ab.Mp3Files = append(c.ab.Mp3Files, dto.Mp3File{Number: i, FileName: localFileName, Size: iaFile.Size, Duration: iaFile.Length})
		jd.AddJob(i, ia.DownloadFile, c.ab.OutputDir, localFileName, item.Server, item.Dir, iaFile.Name, i, iaFile.Size, c.updateFileProgress)
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

	item := c.ab.IAItem
	for !c.stopFlag && percent <= 100 {
		var totalSize = item.TotalSize
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
			totalBytesDownloaded = item.TotalSize
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

			elapsedH := utils.SecondsToTime(elapsed)
			bytesH := utils.BytesToHuman(totalBytesDownloaded)
			filesH := fmt.Sprintf("%d/%d", filesDownloaded, len(item.Files))
			speedH := utils.SpeedToHuman(speed)
			etaH := utils.SecondsToTime(eta)

			c.mq.SendMessage(mq.DownloadController, mq.DownloadPage, &dto.DownloadProgress{Elapsed: elapsedH, Percent: percent, Files: filesH, Bytes: bytesH, Speed: speedH, ETA: etaH}, false)
		}
		time.Sleep(mq.PullFrequency)
	}
}
