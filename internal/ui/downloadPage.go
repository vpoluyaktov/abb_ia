package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rivo/tview"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/dto"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/mq"
)

type DownloadPage struct {
	mq            *mq.Dispatcher
	grid          *tview.Grid
	infoPanel     *infoPanel
	downloadTable *table
}

func newDownloadPage(dispatcher *mq.Dispatcher) *DownloadPage {
	p := &DownloadPage{}
	p.mq = dispatcher
	p.mq.RegisterListener(mq.DownloadPage, p.dispatchMessage)

	p.grid = tview.NewGrid()
	p.grid.SetRows(7, -1, 7)
	p.grid.SetColumns(0)

	// book info section
	infoSection := tview.NewGrid()
	infoSection.SetColumns(-2, -1)
	infoSection.SetBorder(true)
	infoSection.SetTitle(" Audiobook information: ")
	infoSection.SetTitleAlign(tview.AlignLeft)
	p.infoPanel = newInfoPanel()
	infoSection.AddItem(p.infoPanel.t, 0, 0, 1, 1, 0, 0, true)
	f := newForm()
	f.SetHorizontal(false)
	f.f.SetButtonsAlign(tview.AlignRight)
	f.AddButton("Stop", p.stopConfirmation)
	infoSection.AddItem(f.f, 0, 1, 1, 1, 0, 0, false)
	p.grid.AddItem(infoSection, 0, 0, 1, 1, 0, 0, false)

	// files downnload section
	downloadSection := tview.NewGrid()
	downloadSection.SetColumns(-1)
	downloadSection.SetTitle(" Downloading items... ")
	downloadSection.SetTitleAlign(tview.AlignLeft)
	downloadSection.SetBorder(true)

	p.downloadTable = newTable()
	p.downloadTable.setHeaders("  # ", "File name", "Format", "Duration", "Total Size", "Download progress")
	p.downloadTable.setWeights(1, 2, 1, 1, 1, 5)
	p.downloadTable.setAlign(tview.AlignRight, tview.AlignLeft, tview.AlignLeft, tview.AlignRight, tview.AlignRight, tview.AlignLeft)
	// p.downloadTable.t.SetSelectionChangedFunc(p.updateDetails)
	downloadSection.AddItem(p.downloadTable.t, 0, 0, 1, 1, 0, 0, true)
	p.grid.AddItem(downloadSection, 1, 0, 1, 1, 0, 0, true)

	// download progress section
	progressSection := tview.NewGrid()
	progressSection.SetColumns(-2, -1)
	progressSection.SetBorder(true)
	progressSection.SetTitle(" Download progress: ")
	progressSection.SetTitleAlign(tview.AlignLeft)
	p.grid.AddItem(progressSection, 2, 0, 1, 1, 0, 0, false)

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
	case *dto.DownloadProgress:
		p.updateDownloadProgress(dto)
	default:
		m.UnsupportedTypeError(mq.DownloadPage)
	}
}

func (p *DownloadPage) displayBookInfo(ab *dto.Audiobook) {
	p.infoPanel.clear()
	// p.infoPanel.appendRow("", "")
	p.infoPanel.appendRow("Title:", ab.Title)
	p.infoPanel.appendRow("Author:", ab.Author)
	p.infoPanel.appendRow("Duration:", ab.IAItem.TotalLengthH)
	p.infoPanel.appendRow("Size:", ab.IAItem.TotalSizeH)
	p.infoPanel.appendRow("Files", strconv.Itoa(ab.IAItem.FilesCount))

	p.downloadTable.clear()
	p.downloadTable.showHeader()
	for i, f := range ab.IAItem.Files {
		p.downloadTable.appendRow(" "+strconv.Itoa(i+1)+" ", f.Name, f.Format, f.LengthH, f.SizeH, "")
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

func (p *DownloadPage) updateDownloadProgress(dp *dto.DownloadProgress) {
	col := 5
	w := p.downloadTable.getColumnWidth(col) - 5
	progressText := fmt.Sprintf(" %3d%% ", dp.Percent)
	barWidth := int((float32((w - len(progressText))) * float32(dp.Percent) / 100))
	progressBar := strings.Repeat("━", barWidth) + strings.Repeat(" ", w-len(progressText)-barWidth)
	cell := p.downloadTable.t.GetCell(dp.FileId+1, col)
	cell.SetExpansion(0)
	cell.SetMaxWidth(50)
	cell.Text = fmt.Sprintf("%s [%s]", progressText, progressBar)
	p.mq.SendMessage(mq.DownloadPage, mq.TUI, &dto.DrawCommand{Primitive: nil}, true)
}
