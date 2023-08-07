package controller

import (
	"github.com/vpoluyaktov/abb_ia/internal/dto"
	"github.com/vpoluyaktov/abb_ia/internal/logger"
	"github.com/vpoluyaktov/abb_ia/internal/mq"
)

type ChaptersController struct {
	mq       *mq.Dispatcher
	ab       *dto.Audiobook
	stopFlag bool
}

func NewChaptersController(dispatcher *mq.Dispatcher) *ChaptersController {
	dc := &ChaptersController{}
	dc.mq = dispatcher
	dc.mq.RegisterListener(mq.ChaptersController, dc.dispatchMessage)
	return dc
}

func (c *ChaptersController) checkMQ() {
	m := c.mq.GetMessage(mq.ChaptersController)
	if m != nil {
		c.dispatchMessage(m)
	}
}

func (c *ChaptersController) dispatchMessage(m *mq.Message) {
	switch dto := m.Dto.(type) {
	case *dto.EncodeCommand:
		go c.startChapters(dto)
	case *dto.StopCommand:
		go c.stopChapters(dto)
	default:
		m.UnsupportedTypeError(mq.ChaptersController)
	}
}

func (c *ChaptersController) stopChapters(cmd *dto.StopCommand) {
	c.stopFlag = true
	logger.Debug(mq.ChaptersController + ": Received StopChapters command")
}

func (c *ChaptersController) startChapters(cmd *dto.EncodeCommand) {

	c.ab = cmd.Audiobook

	// c.mq.SendMessage(mq.ChaptersController, mq.Footer, &dto.SetBusyIndicator{Busy: false}, false)
	// c.mq.SendMessage(mq.ChaptersController, mq.Footer, &dto.UpdateStatus{Message: ""}, false)
	// c.mq.SendMessage(mq.ChaptersController, mq.ChaptersPage, &dto.ChaptersComplete{Audiobook: cmd.Audiobook}, true)
}
