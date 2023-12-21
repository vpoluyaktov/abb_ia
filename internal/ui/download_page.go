package ui

import (
	"fmt"
	"strconv"
	"strings"

	"abb_ia/internal/dto"
	"abb_ia/internal/mq"
	"abb_ia/internal/utils"

	"github.com/vpoluyaktov/tview"
)

type DownloadPage struct {
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

func newDownloadPage(dispatcher *mq.Dispatcher) *DownloadPage {
	p := &DownloadPage{}
	p.mq = dispatcher
	p.mq.RegisterListener(mq.DownloadPage, p.dispatchMessage)

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

	// files downnload section
	p.filesSection = newGrid()
	p.filesSection.SetColumns(-1)
	p.filesSection.SetTitle(" Downloading .mp3 files... ")
	p.filesSection.SetTitleAlign(tview.AlignLeft)
	p.filesSection.SetBorder(true)

	p.filesTable = newTable()
	p.filesTable.setHeaders(" # ", "File name", "Format", "Duration", "Total Size", "Download progress")
	p.filesTable.setWeights(1, 2, 1, 1, 1, 5)
	p.filesTable.setAlign(tview.AlignRight, tview.AlignLeft, tview.AlignLeft, tview.AlignRight, tview.AlignRight, tview.AlignLeft)
	p.filesSection.AddItem(p.filesTable.Table, 0, 0, 1, 1, 0, 0, true)
	p.mainGrid.AddItem(p.filesSection.Grid, 1, 0, 1, 1, 0, 0, true)

	// download progress section
	p.progressSection = newGrid()
	p.progressSection.SetColumns(-1)
	p.progressSection.SetBorder(true)
	p.progressSection.SetTitle(" Download progress: ")
	p.progressSection.SetTitleAlign(tview.AlignLeft)
	p.progressTable = newTable()
	p.progressTable.setWeights(1)
	p.progressTable.setAlign(tview.AlignLeft)
	p.progressTable.SetSelectable(false, false)
	p.progressSection.AddItem(p.progressTable.Table, 0, 0, 1, 1, 0, 0, false)
	p.mainGrid.AddItem(p.progressSection.Grid, 2, 0, 1, 1, 0, 0, false)

	// screen navigation order
	p.mainGrid.SetNavigationOrder(
		p.infoPanel.Table,
		p.filesTable,
		p.progressTable,
	)

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
	case *dto.FileDownloadProgress:
		p.updateFileProgress(dto)
	case *dto.TotalDownloadProgress:
		p.updateTotalProgress(dto)
	case *dto.DownloadComplete:
		p.downloadComplete(dto)
	default:
		m.UnsupportedTypeError(mq.DownloadPage)
	}
}

func (p *DownloadPage) displayBookInfo(ab *dto.Audiobook) {
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
		p.filesTable.appendRow(" "+strconv.Itoa(i+1)+" ", f.Name, f.Format, utils.SecondsToTime(f.Length), utils.BytesToHuman(f.Size), "")
	}
	p.filesTable.ScrollToBeginning()
	ui.SetFocus(p.filesTable.Table)
	ui.Draw()
}

func (p *DownloadPage) stopConfirmation() {
	newYesNoDialog(p.mq, "Stop Confirmation", "Are you sure you want to stop the download?", p.filesSection.Grid, p.stopDownload, func() {})
}

func (p *DownloadPage) stopDownload() {
	p.mq.SendMessage(mq.DownloadPage, mq.DownloadController, &dto.StopCommand{Process: "Download", Reason: "User request"}, false)
	p.mq.SendMessage(mq.DownloadPage, mq.CleanupController, &dto.CleanupCommand{Audiobook: p.ab}, true)
	p.mq.SendMessage(mq.DownloadPage, mq.Frame, &dto.SwitchToPageCommand{Name: "SearchPage"}, false)
}

func (p *DownloadPage) updateFileProgress(dp *dto.FileDownloadProgress) {
	col := 5
	w := p.filesTable.GetColumnWidth(col) - 5
	if w > 0 {
		progressText := fmt.Sprintf(" %3d%% ", dp.Percent)
		barWidth := int((float32((w - len(progressText))) * float32(dp.Percent) / 100))
		if barWidth < 0 {
			barWidth = 0
		}
		fillerWidth := w - len(progressText) - barWidth
		if fillerWidth < 0 {
			fillerWidth = 0
		}
		progressBar := strings.Repeat("━", barWidth) + strings.Repeat(" ", fillerWidth)
		cell := p.filesTable.GetCell(dp.FileId+1, col)
		cell.SetExpansion(p.filesTable.colWeight[col])
		cell.Text = fmt.Sprintf("%s |%s|", progressText, progressBar)
		ui.Draw()
	}
}

func (p *DownloadPage) updateTotalProgress(dp *dto.TotalDownloadProgress) {
	if p.progressTable.GetRowCount() == 0 {
		for i := 0; i < 2; i++ {
			p.progressTable.appendRow("")
		}
	}
	infoCell := p.progressTable.GetCell(0, 0)
	progressCell := p.progressTable.GetCell(1, 0)
	infoCell.Text = fmt.Sprintf("  [yellow]Time elapsed: [white]%10s | [yellow]Downloaded: [white]%10s | [yellow]Files: [white]%10s | [yellow]Speed: [white]%12s | [yellow]ETA: [white]%10s", dp.Elapsed, dp.Bytes, dp.Files, dp.Speed, dp.ETA)

	col := 0
	w := p.progressTable.GetColumnWidth(col) - 5
	if w > 0 {
		progressText := fmt.Sprintf(" %3d%% ", dp.Percent)
		barWidth := int((float32((w - len(progressText))) * float32(dp.Percent) / 100))
		progressBar := strings.Repeat("▒", barWidth) + strings.Repeat(" ", w-len(progressText)-barWidth)
		progressCell.Text = fmt.Sprintf("%s |%s|", progressText, progressBar)
		ui.Draw()
	}
}

func (p *DownloadPage) downloadComplete(c *dto.DownloadComplete) {
	ab := c.Audiobook
	if ab.Config.IsReEncodeFiles() {
		p.mq.SendMessage(mq.DownloadPage, mq.EncodingController, &dto.EncodeCommand{Audiobook: c.Audiobook}, true)
		p.mq.SendMessage(mq.DownloadPage, mq.Frame, &dto.SwitchToPageCommand{Name: "EncodingPage"}, false)
	} else {
		p.mq.SendMessage(mq.DownloadPage, mq.ChaptersPage, &dto.DisplayBookInfoCommand{Audiobook: c.Audiobook}, true)
		p.mq.SendMessage(mq.DownloadPage, mq.ChaptersController, &dto.ChaptersCreate{Audiobook: c.Audiobook}, true)
		p.mq.SendMessage(mq.DownloadPage, mq.Frame, &dto.SwitchToPageCommand{Name: "ChaptersPage"}, false)
	}
}
