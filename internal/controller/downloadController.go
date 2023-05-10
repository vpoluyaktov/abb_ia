package controller

import (
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/dto"
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
	switch t := m.Type; t {
	case dto.StopCommandType:
		if cmd, ok := m.Dto.(*dto.StopCommand); ok {
			go c.stopDownload(cmd)
		} else {
			m.DtoCastError(mq.DownloadController)
		}

	default:
		m.UnsupportedTypeError(mq.DownloadController)
	}
}

func (c *DownloadController) stopDownload(cmd *dto.StopCommand) {

}