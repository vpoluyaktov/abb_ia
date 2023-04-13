package ui

import (
	"code.rocketnine.space/tslocum/cview"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/mq"
)

type footer struct {
	view       *cview.TextView
	dispatcher *mq.Dispatcher
}

func newFooter(dispatcher *mq.Dispatcher) *footer {
	f := &footer{}
	f.dispatcher = dispatcher
	f.view = cview.NewTextView()
	f.view.SetText("v.0.0.1")
	f.view.SetTextAlign(cview.AlignRight)
	f.view.SetBorder(false)
	f.view.SetTextColor(footerFgColor)
	f.view.SetBackgroundColor(footerBgColor)
	return f
}
