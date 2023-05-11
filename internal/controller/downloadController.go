package controller

import (
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/dto"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/logger"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/mq"
)

type DownloadController struct {
	mq *mq.Dispatcher
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
	logger.Debug(mq.DownloadController + ": Received StopDownload command")
}

func (c *DownloadController) startDownload(cmd *dto.DownloadCommand) {
	logger.Debug(mq.DownloadController + ": Received StartDownload command with IA item: " + cmd.String())
}
