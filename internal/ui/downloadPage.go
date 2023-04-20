package ui

import (
	"github.com/rivo/tview"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/dto"

	"github.com/vpoluyaktov/audiobook_creator_IA/internal/mq"
)

type DownloadPage struct {
	mq   *mq.Dispatcher
	grid *tview.Grid
	downloadSection *tview.Grid
	downloadTable *table
}

func newDownloadPage(dispatcher *mq.Dispatcher) *DownloadPage {
	p := &DownloadPage{}
	p.mq = dispatcher
	p.mq.RegisterListener(mq.DownloadPage, p.dispatchMessage)

	p.grid = tview.NewGrid()
	p.grid.SetRows(5, -1, -1)
	p.grid.SetColumns(0)

	// Download section
	p.downloadSection = tview.NewGrid()
	p.downloadSection.SetColumns(-1)
	p.downloadSection.SetTitle(" Downloading items...")
	p.downloadSection.SetTitleAlign(tview.AlignLeft)
	p.downloadSection.SetBorder(true)

	p.downloadTable = newTable()
	p.downloadTable.setHeaders("Author", "Title", "Files", "Duration (HH:MM:SS)", "Total Size")
	p.downloadTable.setWidths(3, 6, 2, 1, 1)
	p.downloadTable.setAlign(tview.AlignLeft, tview.AlignLeft, tview.AlignRight, tview.AlignRight, tview.AlignRight)
	// p.downloadTable.t.SetSelectionChangedFunc(p.updateDetails)
	p.downloadSection.AddItem(p.downloadTable.t, 0, 0, 1, 1, 0, 0, true)
	p.grid.AddItem(p.downloadSection, 1, 0, 1, 1, 0, 0, true)

	return p
}

func (p *DownloadPage) checkMQ() {
	m := p.mq.GetMessage(mq.DownloadPage)
	if m != nil {
		p.dispatchMessage(m)
	}
}

func (p *DownloadPage) dispatchMessage(m *mq.Message) {
	switch t := m.Type; t {
	case dto.IAItemType:
		if r, ok := m.Dto.(*dto.IAItem); ok {
			go p.updateResult(r)
		} else {
			m.DtoCastError(mq.DownloadPage)
		}

	default:
		m.UnsupportedTypeError(mq.DownloadPage)
	}
}

func (p *DownloadPage) updateResult(i *dto.IAItem) {

}
