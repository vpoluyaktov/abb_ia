package ui

import (
	"github.com/rivo/tview"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/dto"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/mq"
)

type footer struct {
	mq      *mq.Dispatcher
	grid    *tview.Grid
	busy    *tview.TextView
	status  *tview.TextView
	version *tview.TextView
}

func newFooter(dispatcher *mq.Dispatcher) *footer {
	f := &footer{}
	f.mq = dispatcher
	f.mq.RegisterListener(mq.Footer, f.dispatchMessage)

	f.busy = tview.NewTextView()
	f.busy.SetText("")
	f.busy.SetTextAlign(tview.AlignCenter)
	f.busy.SetBorder(false)
	f.busy.SetTextColor(footerFgColor)
	f.busy.SetBackgroundColor(footerBgColor)

	f.status = tview.NewTextView()
	f.status.SetText("")
	f.status.SetTextAlign(tview.AlignLeft)
	f.status.SetBorder(false)
	f.status.SetTextColor(footerFgColor)
	f.status.SetBackgroundColor(footerBgColor)

	f.version = tview.NewTextView()
	f.version.SetText("v.0.0.1")
	f.version.SetTextAlign(tview.AlignRight)
	f.version.SetBorder(false)
	f.version.SetTextColor(footerFgColor)
	f.version.SetBackgroundColor(footerBgColor)

	f.grid = tview.NewGrid()
	f.grid.SetColumns(5, -1, 10)
	f.grid.AddItem(f.busy, 0, 0, 1, 1, 0, 0, false)
	f.grid.AddItem(f.status, 0, 1, 1, 1, 0, 0, false)
	f.grid.AddItem(f.version, 0, 2, 1, 1, 0, 0, false)

	return f
}

func (f *footer) checkMQ() {
	m := f.mq.GetMessage(mq.Footer)
	if m != nil {
		f.dispatchMessage(m)
	}
}

func (f *footer) dispatchMessage(m *mq.Message) {
	switch t := m.Type; t {
	case dto.UpdateStatusType:
		if s, ok := m.Dto.(*dto.UpdateStatus); ok {
			f.updateStatus(s)
		} else {
			m.DtoCastError(mq.Footer)
		}
	case dto.SetBusyIndicatorType:
		if c, ok := m.Dto.(*dto.SetBusyIndicator); ok {
			f.updateBusyIndicator(c)
		} else {
			m.DtoCastError(mq.Footer)
		}	
		

	default:
		m.UnsupportedTypeError(mq.Footer)
	}
}

func (f *footer) updateStatus(s *dto.UpdateStatus) {
	f.status.SetText(s.Message)
}

func (f *footer) updateBusyIndicator(c *dto.SetBusyIndicator) {
	if c.Busy {
		f.busy.SetText("XXX")
	} else {
		f.busy.SetText("")
	}
}
