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

	f.dispatcher.RegisterListener(mq.Frame, f.dispatchMessage)
	return f
}

func (f *frame) addHeader(header *header) {
	f.grid.AddItem(header.view, 0, 1, 1, 1, 0, 0, false)
}

func (f *frame) addFooter(footer *footer) {
	f.grid.AddItem(footer.view, 2, 1, 1, 1, 0, 0, false)
}

func (f *frame) addPage(name string, g *tview.Grid) {
	if f.pages == nil {
		f.pages = tview.NewPages()
		f.grid.AddItem(f.pages, 1, 1, 1, 1, 0, 0, false)
	}
	f.pages.AddPage(name, g, true, true)
}

func (f *frame) removePage(name string) {
	f.pages.RemovePage(name)
}

func (f *frame) showPage(name string) {
	f.pages.SendToFront(name)
}

func (f *frame) checkMQ() {
	m := f.dispatcher.GetMessage(mq.Frame)
	if m != nil {
		f.dispatchMessage(m)
	}
}

func (f *frame) dispatchMessage(m *mq.Message) {
	switch t := m.Type; t {
	case dto.AddPageCommandType:
		if c, ok := m.Dto.(*dto.AddPageCommand); ok {
			f.addPage(c.Name, c.Grid)
		} else {
			m.DtoCastError()
		}
	case dto.RemovePageCommandType:
		if c, ok := m.Dto.(*dto.RemovePageCommand); ok {
			f.removePage(c.Name)
		} else {
			m.DtoCastError()
		}
	case dto.ShowPageCommandType:
		if c, ok := m.Dto.(*dto.ShowPageCommand); ok {
			f.showPage(c.Name)
		} else {
			m.DtoCastError()
		}

	default:
		m.UnsupportedTypeError()
	}
}
