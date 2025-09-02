package ui

import (
	"fmt"

	"abb_ia/internal/dto"
	"abb_ia/internal/logger"
	"abb_ia/internal/mq"

	"github.com/vpoluyaktov/tview"
)

type header struct {
	mq      *mq.Dispatcher
	grid    *tview.Grid
	appName *tview.TextView
	version *tview.TextView
}

func newHeader(dispatcher *mq.Dispatcher) *header {
	h := &header{}
	h.mq = dispatcher
	h.mq.RegisterListener(mq.Header, h.dispatchMessage)

	h.appName = tview.NewTextView()
	h.appName.SetText(" Audiobook Builder")
	h.appName.SetBorder(false)
	h.appName.SetTextColor(headerFgColor)
	h.appName.SetBackgroundColor(headerBGColor)

	h.version = tview.NewTextView()
	h.version.SetText("Internet Archive version ")
	h.version.SetTextAlign(tview.AlignRight)
	h.version.SetBorder(false)
	h.version.SetTextColor(footerFgColor)
	h.version.SetBackgroundColor(footerBgColor)

	h.grid = tview.NewGrid()
	h.grid.SetColumns(-1, 50)
	h.grid.AddItem(h.appName, 0, 0, 1, 1, 0, 0, false)
	h.grid.AddItem(h.version, 0, 1, 1, 1, 0, 0, false)

	return h
}

func (h *header) checkMQ() {
	m, err := h.mq.GetMessage(mq.Header)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to get message for Header: %v", err))
		return
	}
	if m != nil {
		h.dispatchMessage(m)
	}
}

func (h *header) dispatchMessage(m *mq.Message) {
	switch dto := m.Dto.(type) {
	case *dto.DrawCommand:
		if dto.Primitive == nil {
			return
		}
	default:
		m.UnsupportedTypeError(mq.Header)
	}
}
