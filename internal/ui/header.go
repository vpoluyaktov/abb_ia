package ui

import (
	"code.rocketnine.space/tslocum/cview"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/mq"
)

type header struct {
	view   *cview.TextView
	dispatcher *mq.Dispatcher
}

func newHeader(dispatcher *mq.Dispatcher) *header {
	h := &header{}
	h.dispatcher = dispatcher
	h.view = cview.NewTextView()
	h.view.SetText("Audiobook Creator (Internet Archive version)")
	h.view.SetBorder(false)
	h.view.SetTextColor(headerFgColor)
	h.view.SetBackgroundColor(headerBGColor)
	return h
}
