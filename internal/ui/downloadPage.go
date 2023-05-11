package ui

import (
	"github.com/rivo/tview"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/dto"

	"github.com/vpoluyaktov/audiobook_creator_IA/internal/mq"
)

type DownloadPage struct {
	mq              *mq.Dispatcher
	grid            *tview.Grid
	infoSection     *tview.Grid
	downloadSection *tview.Grid
	downloadTable   *table
}

func newDownloadPage(dispatcher *mq.Dispatcher) *DownloadPage {
	p := &DownloadPage{}
	p.mq = dispatcher
	p.mq.RegisterListener(mq.DownloadPage, p.dispatchMessage)

	p.grid = tview.NewGrid()
	p.grid.SetRows(10, -1)
	p.grid.SetColumns(0)

	// information section
	p.infoSection = tview.NewGrid()
	p.infoSection.SetColumns(-2, -1)
	p.infoSection.SetBorder(true)
	p.infoSection.SetTitle(" Audiobook information: ")
	p.infoSection.SetTitleAlign(tview.AlignLeft)
	f := newForm()
	f.SetHorizontal(false)
	f.AddInputField("Search criteria", "", 40, nil, func(t string) {})
	p.infoSection.AddItem(f.f, 0, 0, 1, 1, 0, 0, true)
	f = newForm()
	f.SetHorizontal(false)
	f.f.SetButtonsAlign(tview.AlignRight)
	f.AddButton("Stop", p.stopConfirmation)
	p.infoSection.AddItem(f.f, 0, 1, 1, 1, 0, 0, false)
	p.grid.AddItem(p.infoSection, 0, 0, 1, 1, 0, 0, false)

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
	switch dto := m.Dto.(type) {
	case *dto.IAItem:
		go p.updateResult(dto)
	default:
		m.UnsupportedTypeError(mq.DownloadPage)
	}
}

func (p *DownloadPage) updateResult(i *dto.IAItem) {

}

func (p *DownloadPage) stopConfirmation() {
	newYesNoDialog(p.mq, "Stop Confirmation", "Are you sure you want to stop the download?", p.stopDownload, func() {})
}

func (p *DownloadPage) stopDownload() {
	// Stop the download here
	p.mq.SendMessage(mq.DownloadPage, mq.DownloadController, &dto.StopCommand{Process: "Download"}, false)
	p.mq.SendMessage(mq.DownloadPage, mq.Frame, &dto.SwitchToPageCommand{Name: "SearchPage"}, false)
}
