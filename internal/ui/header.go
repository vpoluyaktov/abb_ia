package ui

import (
	"github.com/rivo/tview"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/dto"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/mq"
)

type header struct {
	view *tview.TextView
	mq   *mq.Dispatcher
}

func newHeader(dispatcher *mq.Dispatcher) *header {
	h := &header{}
	h.mq = dispatcher
	h.view = tview.NewTextView()
	h.view.SetText("Audiobook Creator (Internet Archive version)")
	h.view.SetBorder(false)
	h.view.SetTextColor(headerFgColor)
	h.view.SetBackgroundColor(headerBGColor)

	h.mq.RegisterListener(mq.Header, h.dispatchMessage)
	return h
}

func (h *header) checkMQ() {
	m := h.mq.GetMessage(mq.Header)
	if m != nil {
		h.dispatchMessage(m)
	}
}

func (h *header) dispatchMessage(m *mq.Message) {
	switch dto := m.Dto.(type) {
	case *dto.DrawCommand:
			if dto.Primitive == nil {
			}
	default:
		m.UnsupportedTypeError(mq.Header)
	}
}
