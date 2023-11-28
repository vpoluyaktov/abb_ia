package ui

import (
	"fmt"
	"strconv"
	"strings"

	"abb_ia/internal/dto"
	"abb_ia/internal/mq"
	"abb_ia/internal/utils"

	"github.com/rivo/tview"
)

type EncodingPage struct {
	mq              *mq.Dispatcher
	mainGrid        *grid
	infoSection     *grid
	infoPanel       *infoPanel
	filesSection    *grid
	filesTable      *table
	progressSection *grid
	progressTable   *table
	ab              *dto.Audiobook
}

func newEncodingPage(dispatcher *mq.Dispatcher) *EncodingPage {
	p := &EncodingPage{}
	p.mq = dispatcher
	p.mq.RegisterListener(mq.EncodingPage, p.dispatchMessage)

	p.mainGrid = newGrid()
	p.mainGrid.SetRows(7, -1, 4)
	p.mainGrid.SetColumns(0)

	// book info section
	p.infoSection = newGrid()
	p.infoSection.SetColumns(-2, -1)
	p.infoSection.SetBorder(true)
	p.infoSection.SetTitle(" Audiobook information: ")
	p.infoSection.SetTitleAlign(tview.AlignLeft)
	p.infoPanel = newInfoPanel()
	p.infoSection.AddItem(p.infoPanel.Table, 0, 0, 1, 1, 0, 0, true)
	f := newForm()
	f.SetHorizontal(false)
	f.SetButtonsAlign(tview.AlignRight)
	f.AddButton("Stop", p.stopConfirmation)
	p.infoSection.AddItem(f.Form, 0, 1, 1, 1, 0, 0, false)
	p.mainGrid.AddItem(p.infoSection.Grid, 0, 0, 1, 1, 0, 0, false)

	// files re-encoding section
	p.filesSection = newGrid()
	p.filesSection.SetColumns(-1)
	p.filesSection.SetTitle(" Re-encodinging .mp3 files to the same bitrate... ")
	p.filesSection.SetTitleAlign(tview.AlignLeft)
	p.filesSection.SetBorder(true)

	p.filesTable = newTable()
	p.filesTable.setHeaders("  # ", "File name", "Format", "Duration", "Total Size", "Encoding progress")
	p.filesTable.setWeights(1, 2, 1, 1, 1, 5)
	p.filesTable.setAlign(tview.AlignRight, tview.AlignLeft, tview.AlignLeft, tview.AlignRight, tview.AlignRight, tview.AlignLeft)
	p.filesSection.AddItem(p.filesTable.Table, 0, 0, 1, 1, 0, 0, true)
	p.mainGrid.AddItem(p.filesSection.Grid, 1, 0, 1, 1, 0, 0, true)

	// encoding progress section
	p.progressSection = newGrid()
	p.progressSection.SetColumns(-1)
	p.progressSection.SetBorder(true)
	p.progressSection.SetTitle(" Encoding progress: ")
	p.progressSection.SetTitleAlign(tview.AlignLeft)
	p.progressTable = newTable()
	p.progressTable.setWeights(1)
	p.progressTable.setAlign(tview.AlignLeft)
	p.progressTable.SetSelectable(false, false)
	p.progressSection.AddItem(p.progressTable.Table, 0, 0, 1, 1, 0, 0, false)
	p.mainGrid.AddItem(p.progressSection.Grid, 2, 0, 1, 1, 0, 0, false)

	p.mainGrid.SetNavigationOrder(
		p.infoPanel.Table,
		p.filesTable,
		p.progressTable,
	)

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
	p.ab = ab
	p.infoPanel.clear()
	// p.infoPanel.appendRow("", "")
	p.infoPanel.appendRow("Title:", ab.Title)
	p.infoPanel.appendRow("Author:", ab.Author)
	p.infoPanel.appendRow("Duration:", utils.SecondsToTime(ab.IAItem.TotalLength))
	p.infoPanel.appendRow("Size:", utils.BytesToHuman(ab.IAItem.TotalSize))
	p.infoPanel.appendRow("Files", strconv.Itoa(len(ab.IAItem.AudioFiles)))

	p.filesTable.Clear()
	p.filesTable.showHeader()
	for i, f := range ab.IAItem.AudioFiles {
		p.filesTable.appendRow(" "+strconv.Itoa(i+1)+" ", f.Name, fmt.Sprintf("MP3 %d kb/s", ab.Config.GetBitRate()), utils.SecondsToTime(f.Length), utils.BytesToHuman(f.Size), "")
	}
	p.filesTable.ScrollToBeginning()
	ui.SetFocus(p.filesTable.Table)
	ui.Draw()
}

func (p *EncodingPage) stopConfirmation() {
	newYesNoDialog(p.mq, "Stop Confirmation", "Are you sure you want to stop encoding?", p.filesSection.Grid, p.stopEncoding, func() {})
}

func (p *EncodingPage) stopEncoding() {
	p.mq.SendMessage(mq.EncodingPage, mq.EncodingController, &dto.StopCommand{Process: "Encoding", Reason: "User request"}, true)
	p.mq.SendMessage(mq.EncodingPage, mq.CleanupController, &dto.CleanupCommand{Audiobook: p.ab}, true)
	p.mq.SendMessage(mq.EncodingPage, mq.Frame, &dto.SwitchToPageCommand{Name: "SearchPage"}, true)
}

func (p *EncodingPage) updateFileProgress(dp *dto.EncodingFileProgress) {
	col := 5
	w := p.filesTable.getColumnWidth(col) - 4
	progressText := fmt.Sprintf(" %3d%% ", dp.Percent)
	barWidth := int((float32((w - len(progressText))) * float32(dp.Percent) / 100))
	progressBar := strings.Repeat("━", barWidth) + strings.Repeat(" ", w-len(progressText)-barWidth)
	cell := p.filesTable.GetCell(dp.FileId+1, col)
	cell.SetExpansion(0)
	cell.SetMaxWidth(50)
	cell.Text = fmt.Sprintf("%s |%s|", progressText, progressBar)
	// p.encodingTable.t.Select(dp.FileId+1, col)
	ui.Draw()
}

func (p *EncodingPage) updateTotalProgress(dp *dto.EncodingProgress) {
	if p.progressTable.GetRowCount() == 0 {
		for i := 0; i < 2; i++ {
			p.progressTable.appendRow("")
		}
	}
	infoCell := p.progressTable.GetCell(0, 0)
	progressCell := p.progressTable.GetCell(1, 0)
	infoCell.Text = fmt.Sprintf("  [yellow]Time elapsed: [white]%10s | [yellow]Files: [white]%10s | [yellow]Speed: [white]%10s | [yellow]ETA: [white]%10s", dp.Elapsed, dp.Files, dp.Speed, dp.ETA)

	col := 0
	w := p.progressTable.getColumnWidth(col) - 5
	progressText := fmt.Sprintf(" %3d%% ", dp.Percent)
	barWidth := int((float32((w - len(progressText))) * float32(dp.Percent) / 100))
	progressBar := strings.Repeat("▒", barWidth) + strings.Repeat(" ", w-len(progressText)-barWidth)
	// progressCell.SetExpansion(0)
	// progressCell.SetMaxWidth(0)
	progressCell.Text = fmt.Sprintf("%s |%s|", progressText, progressBar)
	ui.Draw()
}

func (p *EncodingPage) encodingComplete(c *dto.EncodingComplete) {
	p.mq.SendMessage(mq.EncodingPage, mq.ChaptersController, &dto.ChaptersCreate{Audiobook: c.Audiobook}, true)
	p.mq.SendMessage(mq.EncodingPage, mq.Frame, &dto.SwitchToPageCommand{Name: "ChaptersPage"}, false)
}
