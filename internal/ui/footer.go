package ui

import (
	"github.com/rivo/tview"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/mq"
)

type footer struct {
	view       *tview.TextView
	dispatcher *mq.Dispatcher
}

func newFooter(dispatcher *mq.Dispatcher) *footer {
	f := &footer{}
	f.dispatcher = dispatcher
	f.view = tview.NewTextView()
	f.view.SetText("v.0.0.1")
	f.view.SetTextAlign(tview.AlignRight)
	f.view.SetBorder(false)
	f.view.SetTextColor(footerFgColor)
	f.view.SetBackgroundColor(footerBgColor)
	return f
}
