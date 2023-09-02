package controller

import (
	"path/filepath"

	"github.com/vpoluyaktov/abb_ia/internal/dto"
	"github.com/vpoluyaktov/abb_ia/internal/ffmpeg"
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
	case *dto.ChaptersCreate:
		go c.createChapters(dto)
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

func (c *ChaptersController) createChapters(cmd *dto.ChaptersCreate) {

	logger.Debug(mq.ChaptersController + ": Received Create Chapters command")
	c.mq.SendMessage(mq.ChaptersController, mq.Footer, &dto.UpdateStatus{Message: "Calculating book parts and chapters..."}, false)
	c.mq.SendMessage(mq.ChaptersController, mq.Footer, &dto.SetBusyIndicator{Busy: true}, false)

	c.ab = cmd.Audiobook
	c.ab.Chapters = nil

	var offset float64 = 0
	for i, file := range c.ab.IAItem.Files {
		p, _ := ffmpeg.NewFFProbe(filepath.Join("output", c.ab.IAItem.ID, c.ab.IAItem.Dir, file.Name))
		chapter := dto.Chapter{Number: i + 1, Start: offset, End: offset + p.GetDuration(), Duration: p.GetDuration(), Name: p.GetTitle()}
		c.ab.Chapters = append(c.ab.Chapters, chapter)
		offset += p.GetDuration()
		c.mq.SendMessage(mq.ChaptersController, mq.ChaptersPage, &dto.AddChapterCommand{Chapter: &chapter}, true)
	}
	c.mq.SendMessage(mq.ChaptersController, mq.Footer, &dto.SetBusyIndicator{Busy: false}, false)
	c.mq.SendMessage(mq.ChaptersController, mq.Footer, &dto.UpdateStatus{Message: ""}, false)
	c.mq.SendMessage(mq.ChaptersController, mq.ChaptersPage, &dto.ChaptersReady{Audiobook: cmd.Audiobook}, true)
}
