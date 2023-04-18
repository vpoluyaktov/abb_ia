package ui

import (
	"github.com/rivo/tview"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/dto"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/mq"
)

type frame struct {
	dispatcher *mq.Dispatcher
	grid       *tview.Grid
	pages      *tview.Pages
}

func newFrame(dispatcher *mq.Dispatcher) *frame {
	f := &frame{}
	f.dispatcher = dispatcher
	f.grid = tview.NewGrid()
	f.grid.SetRows(1, 0, 1)
	f.grid.SetColumns(1, 0, 1)
	// f.grid.SetBackgroundColor(true)

	f.dispatcher.RegisterListener(mq.Frame, f.dispatchMessage)
	return f
}

func (f *frame) addHeader(header *header) {
	f.grid.AddItem(header.view, 0, 1, 1, 1, 0, 0, false)
}

func (f *frame) addFooter(footer *footer) {
	f.grid.AddItem(footer.view, 2, 1, 1, 1, 0, 0, false)
}

func (f *frame) addPannel(name string, g *tview.Grid) {
	if f.pages == nil {
		f.pages = tview.NewPages()
		f.grid.AddItem(f.pages, 1, 1, 1, 1, 0, 0, false)
	}
	f.pages.AddPage(name, g, true, true)
}

func (f *frame) removePannel(name string) {
	f.pages.RemovePage(name)
}

func (f *frame) showPanel(name string) {
	f.pages.SendToFront(name)
	// f.pages.SetPage(name)
}

func (f *frame) checkMQ() {
	m := f.dispatcher.GetMessage(mq.Frame)
	if m != nil {
		f.dispatchMessage(m)
	}
}

func (f *frame) dispatchMessage(m *mq.Message) {
	switch t := m.Type; t {
	case dto.AddPanelCommandType:
		if c, ok := m.Dto.(*dto.AddPanelCommand); ok {
			f.addPannel(c.Name, c.Grid)
		} else {
			m.DtoCastError()
		}
	case dto.RemovePanelCommandType:
		if c, ok := m.Dto.(*dto.RemovePanelCommand); ok {
			f.removePannel(c.Name)
		} else {
			m.DtoCastError()
		}
	case dto.ShowPanelCommandType:
		if c, ok := m.Dto.(*dto.ShowPanelCommand); ok {
			f.showPanel(c.Name)
		} else {
			m.DtoCastError()
		}

	default:
		m.UnsupportedTypeError()
	}
}
