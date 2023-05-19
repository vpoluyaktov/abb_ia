package controller

import (
	"path/filepath"
	"strconv"

	"github.com/vpoluyaktov/audiobook_creator_IA/internal/config"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/dto"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/ia_client"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/logger"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/mq"
)

type DownloadController struct {
	mq       *mq.Dispatcher
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

func (c *DownloadController) startDownload(dto *dto.DownloadCommand) {
	logger.Debug(mq.DownloadController + ": Received StartDownload command with IA item: " + dto.String())
	item := dto.Audiobook.IAItem
	outputDir := filepath.Join("output", item.ID)

	// display DownloadPage initial content


	// download files
	ia := ia_client.New(config.IsUseMock(), config.IsSaveMock())
	c.stopFlag = false
	for _, f := range item.Files {
		if c.stopFlag {
			break
		}
		if false /* config.ParrallelDownload */ {
			go ia.DownloadFile(outputDir, item.Server, item.Dir, f.Name, c.updateProgress)
		} else {
			ia.DownloadFile(outputDir, item.Server, item.Dir, f.Name, c.updateProgress)
		}
	}
}

func (c *DownloadController) updateProgress(fileName string, percent int) {
	logger.Debug("File: " + fileName + " Downloading progress: " + strconv.Itoa(percent))
}
