package ui

import (
	"github.com/rivo/tview"
	"github.com/vpoluyaktov/abb_ia/internal/dto"
	"github.com/vpoluyaktov/abb_ia/internal/mq"
)

type frame struct {
	mq    *mq.Dispatcher
	grid  *tview.Grid
	pages *tview.Pages
}

func newFrame(dispatcher *mq.Dispatcher) *frame {
	f := &frame{}
	f.mq = dispatcher
	f.grid = tview.NewGrid()
	f.grid.SetRows(1, 0, 1)
	// f.grid.SetColumns(1, 0, 1) // extra space on the right and left side
	f.grid.SetColumns(0)

	f.mq.RegisterListener(mq.Frame, f.dispatchMessage)
	return f
}

func (f *frame) addHeader(header *header) {
	f.grid.AddItem(header.view, 0, 0, 1, 1, 0, 0, false)
}

func (f *frame) addFooter(footer *footer) {
	f.grid.AddItem(footer.grid, 2, 0, 1, 1, 0, 0, false)
}

func (f *frame) addPage(name string, g *tview.Grid) {
	if f.pages == nil {
		f.pages = tview.NewPages()
		f.grid.AddItem(f.pages, 1, 0, 1, 1, 0, 0, false)
	}
	f.pages.AddPage(name, g, true, true)
}

func (f *frame) removePage(name string) {
	f.pages.RemovePage(name)
}

func (f *frame) showPage(name string) {
	f.pages.SendToFront(name)
}

func (f *frame) switchToPage(name string) {
	f.pages.SwitchToPage(name)
}

func (f *frame) checkMQ() {
	m := f.mq.GetMessage(mq.Frame)
	if m != nil {
		f.dispatchMessage(m)
	}
}

func (f *frame) dispatchMessage(m *mq.Message) {
	switch dto := m.Dto.(type) {
	case *dto.AddPageCommand:
		f.addPage(dto.Name, dto.Grid)
	case *dto.RemovePageCommand:
		f.removePage(dto.Name)
	case *dto.ShowPageCommand:
		f.showPage(dto.Name)
	case *dto.SwitchToPageCommand:
		f.switchToPage(dto.Name)
	default:
		m.UnsupportedTypeError(mq.Frame)
	}
}
