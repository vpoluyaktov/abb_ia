package ui

import (
	"strconv"

	"github.com/rivo/tview"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/dto"

	"github.com/vpoluyaktov/audiobook_creator_IA/internal/mq"
)

type DownloadPage struct {
	mq              *mq.Dispatcher
	grid            *tview.Grid
	infoSection     *tview.Grid
	infoTable       *infoTable
	downloadSection *tview.Grid
	downloadTable   *table
}

func newDownloadPage(dispatcher *mq.Dispatcher) *DownloadPage {
	p := &DownloadPage{}
	p.mq = dispatcher
	p.mq.RegisterListener(mq.DownloadPage, p.dispatchMessage)

	p.grid = tview.NewGrid()
	p.grid.SetRows(9, -1)
	p.grid.SetColumns(0)

	// information section
	p.infoSection = tview.NewGrid()
	p.infoSection.SetColumns(-2, -1)
	p.infoSection.SetBorder(true)
	p.infoSection.SetTitle(" Audiobook information: ")
	p.infoSection.SetTitleAlign(tview.AlignLeft)
	p.infoTable = newInfoTable()
	p.infoSection.AddItem(p.infoTable.t, 0, 0, 1, 1, 0, 0, true)
	f := newForm()
	f.SetHorizontal(false)
	f.f.SetButtonsAlign(tview.AlignRight)
	f.AddButton("Stop", p.stopConfirmation)
	p.infoSection.AddItem(f.f, 0, 1, 1, 1, 0, 0, false)
	p.grid.AddItem(p.infoSection, 0, 0, 1, 1, 0, 0, false)

	// Download section
	p.downloadSection = tview.NewGrid()
	p.downloadSection.SetColumns(-1)
	p.downloadSection.SetTitle(" Downloading items... ")
	p.downloadSection.SetTitleAlign(tview.AlignLeft)
	p.downloadSection.SetBorder(true)

	p.downloadTable = newTable()
	p.downloadTable.setHeaders(" # ", "File name", "Format", "Duration (HH:MM:SS)", "Total Size", "Download progress")
	p.downloadTable.setWidths(1, 4, 2, 2, 1, 20)
	p.downloadTable.setAlign(tview.AlignRight, tview.AlignLeft, tview.AlignLeft, tview.AlignRight, tview.AlignRight, tview.AlignLeft)
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
	case *dto.DisplayBookInfoCommand:
		p.displayBookInfo(dto.Audiobook)
	default:
		m.UnsupportedTypeError(mq.DownloadPage)
	}
}

func (p *DownloadPage) displayBookInfo(ab *dto.Audiobook) {
	p.infoTable.clear()
	p.infoTable.appendRow("", "")
	p.infoTable.appendRow("Title:", ab.Title)
	p.infoTable.appendRow("Author:", ab.Author)
	p.infoTable.appendRow("Duration:", ab.IAItem.TotalLengthH)
	p.infoTable.appendRow("Size:", ab.IAItem.TotalSizeH)
	p.infoTable.appendRow("Files", strconv.Itoa(ab.IAItem.FilesCount))

	p.downloadTable.clear()
	p.downloadTable.showHeader()
	for i, f := range ab.IAItem.Files {
		p.downloadTable.appendRow(strconv.Itoa(i+1), f.Name, f.Format, f.LengthH, f.SizeH, "")
	}
	p.downloadTable.t.ScrollToBeginning()
	p.mq.SendMessage(mq.DownloadPage, mq.TUI, &dto.SetFocusCommand{Primitive: p.downloadTable.t}, true)
	p.mq.SendMessage(mq.DownloadPage, mq.TUI, &dto.DrawCommand{Primitive: nil}, true)
}

func (p *DownloadPage) stopConfirmation() {
	newYesNoDialog(p.mq, "Stop Confirmation", "Are you sure you want to stop the download?", p.stopDownload, func() {})
}

func (p *DownloadPage) stopDownload() {
	// Stop the download here
	p.mq.SendMessage(mq.DownloadPage, mq.DownloadController, &dto.StopCommand{Process: "Download", Reason: "User request"}, false)
	p.mq.SendMessage(mq.DownloadPage, mq.Frame, &dto.SwitchToPageCommand{Name: "SearchPage"}, false)
}
