package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rivo/tview"
	"github.com/vpoluyaktov/abb_ia/internal/config"
	"github.com/vpoluyaktov/abb_ia/internal/dto"
	"github.com/vpoluyaktov/abb_ia/internal/mq"
)

type EncodingPage struct {
	mq            *mq.Dispatcher
	grid          *tview.Grid
	infoPanel     *infoPanel
	filesTable    *table
	progressTable *table
}

func newEncodingPage(dispatcher *mq.Dispatcher) *EncodingPage {
	p := &EncodingPage{}
	p.mq = dispatcher
	p.mq.RegisterListener(mq.EncodingPage, p.dispatchMessage)

	p.grid = tview.NewGrid()
	p.grid.SetRows(7, -1, 4)
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

	// files re-encoding section
	filesSection := tview.NewGrid()
	filesSection.SetColumns(-1)
	filesSection.SetTitle(" Encodinging .mp3 files... ")
	filesSection.SetTitleAlign(tview.AlignLeft)
	filesSection.SetBorder(true)

	p.filesTable = newTable()
	p.filesTable.setHeaders("  # ", "File name", "Format", "Duration", "Total Size", "Encoding progress")
	p.filesTable.setWeights(1, 2, 1, 1, 1, 5)
	p.filesTable.setAlign(tview.AlignRight, tview.AlignLeft, tview.AlignLeft, tview.AlignRight, tview.AlignRight, tview.AlignLeft)
	filesSection.AddItem(p.filesTable.t, 0, 0, 1, 1, 0, 0, true)
	p.grid.AddItem(filesSection, 1, 0, 1, 1, 0, 0, true)

	// encoding progress section
	progressSection := tview.NewGrid()
	progressSection.SetColumns(-1)
	progressSection.SetBorder(true)
	progressSection.SetTitle(" Encoding progress: ")
	progressSection.SetTitleAlign(tview.AlignLeft)
	p.progressTable = newTable()
	p.progressTable.setWeights(1)
	p.progressTable.setAlign(tview.AlignLeft)
	p.progressTable.t.SetSelectable(false, false)
	progressSection.AddItem(p.progressTable.t, 0, 0, 1, 1, 0, 0, false)
	p.grid.AddItem(progressSection, 2, 0, 1, 1, 0, 0, false)

	return p
}

func (p *EncodingPage) checkMQ() {
	m := p.mq.GetMessage(mq.EncodingPage)
	if m != nil {
		p.dispatchMessage(m)
	}
}

func (p *EncodingPage) dispatchMessage(m *mq.Message) {
	switch dto := m.Dto.(type) {
	case *dto.DisplayBookInfoCommand:
		p.displayBookInfo(dto.Audiobook)
	case *dto.EncodingFileProgress:
		p.updateFileProgress(dto)
	case *dto.EncodingProgress:
		p.updateTotalProgress(dto)
	case *dto.EncodingComplete:
		p.encodingComplete(dto)
	default:
		m.UnsupportedTypeError(mq.EncodingPage)
	}
}

func (p *EncodingPage) displayBookInfo(ab *dto.Audiobook) {
	p.infoPanel.clear()
	// p.infoPanel.appendRow("", "")
	p.infoPanel.appendRow("Title:", ab.Title)
	p.infoPanel.appendRow("Author:", ab.Author)
	p.infoPanel.appendRow("Duration:", ab.IAItem.TotalLengthH)
	p.infoPanel.appendRow("Size:", ab.IAItem.TotalSizeH)
	p.infoPanel.appendRow("Files", strconv.Itoa(ab.IAItem.FilesCount))

	p.filesTable.clear()
	p.filesTable.showHeader()
	for i, f := range ab.IAItem.Files {
		p.filesTable.appendRow(" "+strconv.Itoa(i+1)+" ", f.Name, f.Format, f.LengthH, f.SizeH, "")
	}
	p.filesTable.t.ScrollToBeginning()
	p.mq.SendMessage(mq.EncodingPage, mq.TUI, &dto.SetFocusCommand{Primitive: p.filesTable.t}, true)
	p.mq.SendMessage(mq.EncodingPage, mq.TUI, &dto.DrawCommand{Primitive: nil}, true)
}

func (p *EncodingPage) stopConfirmation() {
	newYesNoDialog(p.mq, "Stop Confirmation", "Are you sure you want to stop encoding?", p.stopEncoding, func() {})
}

func (p *EncodingPage) stopEncoding() {
	// Stop the encoding here
	p.mq.SendMessage(mq.EncodingPage, mq.EncodingController, &dto.StopCommand{Process: "Encoding", Reason: "User request"}, false)
	p.mq.SendMessage(mq.EncodingPage, mq.Frame, &dto.SwitchToPageCommand{Name: "SearchPage"}, false)
}

func (p *EncodingPage) updateFileProgress(dp *dto.EncodingFileProgress) {
	col := 5
	w := p.filesTable.getColumnWidth(col) - 5
	progressText := fmt.Sprintf(" %3d%% ", dp.Percent)
	barWidth := int((float32((w - len(progressText))) * float32(dp.Percent) / 100))
	progressBar := strings.Repeat("━", barWidth) + strings.Repeat(" ", w-len(progressText)-barWidth)
	cell := p.filesTable.t.GetCell(dp.FileId+1, col)
	cell.SetExpansion(0)
	cell.SetMaxWidth(50)
	cell.Text = fmt.Sprintf("%s |%s|", progressText, progressBar)
	// p.encodingTable.t.Select(dp.FileId+1, col)
	p.mq.SendMessage(mq.EncodingPage, mq.TUI, &dto.DrawCommand{Primitive: nil}, true)
}

func (p *EncodingPage) updateTotalProgress(dp *dto.EncodingProgress) {
	if p.progressTable.GetRowCount() == 0 {
		for i := 0; i < 2; i++ {
			p.progressTable.appendRow("")
		}
	}
	infoCell := p.progressTable.t.GetCell(0, 0)
	progressCell := p.progressTable.t.GetCell(1, 0)
	infoCell.Text = fmt.Sprintf("  [yellow]Time elapsed: [white]%10s | [yellow]Files: [white]%10s | [yellow]Speed: [white]%12s | [yellow]ETA: [white]%10s", dp.Elapsed, dp.Files, dp.Speed, dp.ETA)

	col := 0
	w := p.progressTable.getColumnWidth(col) - 5
	progressText := fmt.Sprintf(" %3d%% ", dp.Percent)
	barWidth := int((float32((w - len(progressText))) * float32(dp.Percent) / 100))
	progressBar := strings.Repeat("▒", barWidth) + strings.Repeat(" ", w-len(progressText)-barWidth)
	// progressCell.SetExpansion(0)
	// progressCell.SetMaxWidth(0)
	progressCell.Text = fmt.Sprintf("%s |%s|", progressText, progressBar)
	p.mq.SendMessage(mq.EncodingPage, mq.TUI, &dto.DrawCommand{Primitive: nil}, true)
}

func (p *EncodingPage) encodingComplete(c *dto.EncodingComplete) {
	if config.IsReEncodeFiles() {
	// 	p.mq.SendMessage(mq.DownloadPage, mq.EncodingPage, &dto.DisplayBookInfoCommand{Audiobook: c.Audiobook}, true)
	// 	p.mq.SendMessage(mq.DownloadPage, mq.EncodingController, &dto.EncodeCommand{Audiobook: c.Audiobook}, true)
	// 	p.mq.SendMessage(mq.DownloadPage, mq.Frame, &dto.SwitchToPageCommand{Name: "EncodingPage"}, false)
	// } else {

	}
}