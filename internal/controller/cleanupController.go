package controller

import (
	"os"

	"github.com/vpoluyaktov/abb_ia/internal/dto"
	"github.com/vpoluyaktov/abb_ia/internal/logger"
	"github.com/vpoluyaktov/abb_ia/internal/mq"
)

type CleanupController struct {
	mq *mq.Dispatcher
	ab *dto.Audiobook
}

func NewCleanupController(dispatcher *mq.Dispatcher) *CleanupController {
	c := &CleanupController{}
	c.mq = dispatcher
	c.mq.RegisterListener(mq.CleanupController, c.dispatchMessage)
	return c
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

	if !(c.ab.Config.IsSaveMock() || c.ab.Config.IsUseMock()) {
		os.RemoveAll(c.ab.OutputDir)
	}
	for _, part := range c.ab.Parts {
		os.Remove(part.AACFile)
		os.Remove(part.FListFile)
		os.Remove(part.MetadataFile)
		if c.ab.Config.IsCopyToOutputDir() {
			os.Remove(part.M4BFile)
		}
	}
	os.Remove(c.ab.CoverFile)

	c.mq.SendMessage(mq.CleanupController, mq.BuildPage, &dto.CleanupComplete{Audiobook: cmd.Audiobook}, true)
}
