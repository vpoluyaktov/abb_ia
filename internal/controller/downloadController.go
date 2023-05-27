package controller

import (
	"path/filepath"

	"github.com/vpoluyaktov/audiobook_creator_IA/internal/config"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/dto"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/ia_client"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/logger"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/mq"
)

type DownloadController struct {
	mq       *mq.Dispatcher
	progress []int
	stopFlag bool
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
	logger.Debug(mq.DownloadController + ": Received StartDownload command with IA item: " + cmd.String())
	c.mq.SendMessage(mq.DownloadController, mq.Footer, &dto.UpdateStatus{Message: "Downloading mp3 files..."}, false)
	c.mq.SendMessage(mq.DownloadController, mq.Footer, &dto.SetBusyIndicator{Busy: true}, false)
	item := cmd.Audiobook.IAItem
	outputDir := filepath.Join("output", item.ID)

	// display DownloadPage initial content

	// download files
	ia := ia_client.New(config.IsUseMock(), config.IsSaveMock())
	c.stopFlag = false
	c.progress = make([]int, len(item.Files))
	for i, f := range item.Files {
		if c.stopFlag {
			break
		}
		if config.IsParallelDownload() {
			go ia.DownloadFile(outputDir, item.Server, item.Dir, f.Name, i, c.updateProgress)
		} else {
			ia.DownloadFile(outputDir, item.Server, item.Dir, f.Name, i, c.updateProgress)
		}
	}
	c.mq.SendMessage(mq.DownloadController, mq.Footer, &dto.SetBusyIndicator{Busy: false}, false)
	c.mq.SendMessage(mq.DownloadController, mq.Footer, &dto.UpdateStatus{Message: ""}, false)
}

func (c *DownloadController) updateProgress(fileId int, fileName string, percent int) {
	if c.progress[fileId] != percent {
		// sent a message only if progress changed
		c.mq.SendMessage(mq.DownloadController, mq.DownloadPage, &dto.DownloadProgress{FileId: fileId, FileName: fileName, Percent: percent}, true)
	}
	c.progress[fileId] = percent
}
