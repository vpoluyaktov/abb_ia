package controller

import (
	"os"

	"github.com/vpoluyaktov/abb_ia/internal/config"
	"github.com/vpoluyaktov/abb_ia/internal/dto"
	"github.com/vpoluyaktov/abb_ia/internal/logger"
	"github.com/vpoluyaktov/abb_ia/internal/mq"
)

type CleanupController struct {
	mq *mq.Dispatcher
	ab *dto.Audiobook
}

func NewCleanupController(dispatcher *mq.Dispatcher) *CleanupController {
	dc := &CleanupController{}
	dc.mq = dispatcher
	dc.mq.RegisterListener(mq.CleanupController, dc.dispatchMessage)
	return dc
}

func (c *CleanupController) checkMQ() {
	m := c.mq.GetMessage(mq.CleanupController)
	if m != nil {
		c.dispatchMessage(m)
	}
}

func (c *CleanupController) dispatchMessage(m *mq.Message) {
	switch dto := m.Dto.(type) {
	case *dto.CleanupCommand:
		go c.cleanUp(dto)
	default:
		m.UnsupportedTypeError(mq.CleanupController)
	}
}

func (c *CleanupController) cleanUp(cmd *dto.CleanupCommand) {
	logger.Info(mq.CleanupController + " received " + cmd.String())
	c.ab = cmd.Audiobook

	if !(config.Instance().IsSaveMock() || config.Instance().IsUseMock()) {
		os.RemoveAll(c.ab.OutputDir)
	}
	for _, part := range c.ab.Parts {
		os.Remove(part.AACFile)
		os.Remove(part.FListFile)
		os.Remove(part.MetadataFile)
		if config.Instance().IsCopyToAudiobookshelf() {
			os.Remove(part.M4BFile)
		}
	}
	os.Remove(c.ab.CoverFile)
}
