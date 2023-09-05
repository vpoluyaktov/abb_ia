package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rivo/tview"
	"github.com/vpoluyaktov/abb_ia/internal/dto"
	"github.com/vpoluyaktov/abb_ia/internal/mq"
)

type BuildPage struct {
	mq           *mq.Dispatcher
	grid         *tview.Grid
	infoPanel    *infoPanel
	buildSection *tview.Grid
	buildTable   *table
	copyTable    *table
}

func newBuildPage(dispatcher *mq.Dispatcher) *BuildPage {
	p := &BuildPage{}
	p.mq = dispatcher
	p.mq.RegisterListener(mq.BuildPage, p.dispatchMessage)

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

	// audiobook build section
	p.buildSection = tview.NewGrid()
	p.buildSection.SetColumns(-1)
	p.buildSection.SetTitle(" Building audiobook... ")
	p.buildSection.SetTitleAlign(tview.AlignLeft)
	p.buildSection.SetBorder(true)

	p.buildTable = newTable()
	p.buildTable.setHeaders(" Part # ", "File name", "Duration", "Total Size", "Build progress")
	p.buildTable.setWeights(1, 2, 1, 1, 5)
	p.buildTable.setAlign(tview.AlignRight, tview.AlignLeft, tview.AlignLeft, tview.AlignRight, tview.AlignRight, tview.AlignLeft)
	p.buildSection.AddItem(p.buildTable.t, 0, 0, 1, 1, 0, 0, true)
	p.grid.AddItem(p.buildSection, 1, 0, 1, 1, 0, 0, true)

	// copy section
	copySection := tview.NewGrid()
	copySection.SetColumns(-1)
	copySection.SetBorder(true)
	copySection.SetTitle(" Copy to audioshelf catalog progress: ")
	copySection.SetTitleAlign(tview.AlignLeft)
	p.copyTable = newTable()
	p.copyTable.setHeaders(" Part # ", "File name", "Duration", "Total Size", "Copy progress")
	p.copyTable.setWeights(1, 2, 1, 1, 5)
	p.copyTable.setAlign(tview.AlignRight, tview.AlignLeft, tview.AlignLeft, tview.AlignRight, tview.AlignRight, tview.AlignLeft)
	copySection.AddItem(p.copyTable.t, 0, 0, 1, 1, 0, 0, false)
	p.grid.AddItem(copySection, 2, 0, 1, 1, 0, 0, false)

	return p
}

func (p *BuildPage) checkMQ() {
	m := p.mq.GetMessage(mq.BuildPage)
	if m != nil {
		p.dispatchMessage(m)
	}
}

func (p *BuildPage) dispatchMessage(m *mq.Message) {
	switch dto := m.Dto.(type) {
	case *dto.DisplayBookInfoCommand:
		p.displayBookInfo(dto.Audiobook)
	case *dto.BuildProgress:
		p.updateBuildProgress(dto)
	case *dto.BuildComplete:
		p.buildComplete(dto)
	case *dto.CopyProgress:
		p.updateCopyProgress(dto)
	case *dto.CopyComplete:
		p.copyComplete(dto)
	default:
		m.UnsupportedTypeError(mq.BuildPage)
	}
}

func (p *BuildPage) displayBookInfo(ab *dto.Audiobook) {
	p.infoPanel.clear()
	// p.infoPanel.appendRow("", "")
	p.infoPanel.appendRow("Title:", ab.Title)
	p.infoPanel.appendRow("Author:", ab.Author)
	p.infoPanel.appendRow("Duration:", ab.IAItem.TotalLengthH)
	p.infoPanel.appendRow("Size:", ab.IAItem.TotalSizeH)
	p.infoPanel.appendRow("Files", strconv.Itoa(ab.IAItem.FilesCount))

	p.buildTable.clear()
	p.buildTable.showHeader()
	// for i, f := range ab.IAItem.Files {
	// 	p.buildTable.appendRow(" "+strconv.Itoa(i+1)+" ", f.Name, f.Format, f.LengthH, f.SizeH, "")
	// }
	p.buildTable.t.ScrollToBeginning()
	p.mq.SendMessage(mq.BuildPage, mq.TUI, &dto.SetFocusCommand{Primitive: p.buildTable.t}, true)
	p.mq.SendMessage(mq.BuildPage, mq.TUI, &dto.DrawCommand{Primitive: nil}, true)
}

func (p *BuildPage) stopConfirmation() {
	newYesNoDialog(p.mq, "Stop Confirmation", "Are you sure you want to stop the build?", p.buildSection, p.stopBuild, func() {})
}

func (p *BuildPage) stopBuild() {
	// Stop the build here
	p.mq.SendMessage(mq.BuildPage, mq.BuildController, &dto.StopCommand{Process: "Build", Reason: "User request"}, false)
	p.mq.SendMessage(mq.BuildPage, mq.Frame, &dto.SwitchToPageCommand{Name: "SearchPage"}, false)
}

func (p *BuildPage) updateBuildProgress(dp *dto.BuildProgress) {
	col := 5
	w := p.buildTable.getColumnWidth(col) - 5
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
	cell := p.buildTable.t.GetCell(dp.FileId+1, col)
	cell.SetExpansion(0)
	cell.SetMaxWidth(50)
	cell.Text = fmt.Sprintf("%s |%s|", progressText, progressBar)
	// p.downloadTable.t.Select(dp.FileId+1, col)
	p.mq.SendMessage(mq.DownloadPage, mq.TUI, &dto.DrawCommand{Primitive: nil}, true)
}

func (p *BuildPage) updateCopyProgress(dp *dto.CopyProgress) {
	col := 5
	w := p.copyTable.getColumnWidth(col) - 5
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
	cell := p.copyTable.t.GetCell(dp.FileId+1, col)
	cell.SetExpansion(0)
	cell.SetMaxWidth(50)
	cell.Text = fmt.Sprintf("%s |%s|", progressText, progressBar)
	// p.downloadTable.t.Select(dp.FileId+1, col)
	p.mq.SendMessage(mq.DownloadPage, mq.TUI, &dto.DrawCommand{Primitive: nil}, true)
}

func (p *BuildPage) buildComplete(c *dto.BuildComplete) {
	// if config.IsReEncodeFiles() {
	// 	p.mq.SendMessage(mq.BuildPage, mq.EncodingPage, &dto.DisplayBookInfoCommand{Audiobook: c.Audiobook}, true)
	// 	p.mq.SendMessage(mq.BuildPage, mq.EncodingController, &dto.EncodeCommand{Audiobook: c.Audiobook}, true)
	// 	p.mq.SendMessage(mq.BuildPage, mq.Frame, &dto.SwitchToPageCommand{Name: "EncodingPage"}, false)
	// } else {
	// 	p.mq.SendMessage(mq.BuildPage, mq.ChaptersPage, &dto.DisplayBookInfoCommand{Audiobook: c.Audiobook}, true)
	// 	p.mq.SendMessage(mq.BuildPage, mq.ChaptersController, &dto.ChaptersCreate{Audiobook: c.Audiobook}, true)
	// 	p.mq.SendMessage(mq.BuildPage, mq.Frame, &dto.SwitchToPageCommand{Name: "ChaptersPage"}, false)
	// }
}

func (p *BuildPage) copyComplete(c *dto.CopyComplete) {
	// if config.IsReEncodeFiles() {
	// 	p.mq.SendMessage(mq.BuildPage, mq.EncodingPage, &dto.DisplayBookInfoCommand{Audiobook: c.Audiobook}, true)
	// 	p.mq.SendMessage(mq.BuildPage, mq.EncodingController, &dto.EncodeCommand{Audiobook: c.Audiobook}, true)
	// 	p.mq.SendMessage(mq.BuildPage, mq.Frame, &dto.SwitchToPageCommand{Name: "EncodingPage"}, false)
	// } else {
	// 	p.mq.SendMessage(mq.BuildPage, mq.ChaptersPage, &dto.DisplayBookInfoCommand{Audiobook: c.Audiobook}, true)
	// 	p.mq.SendMessage(mq.BuildPage, mq.ChaptersController, &dto.ChaptersCreate{Audiobook: c.Audiobook}, true)
	// 	p.mq.SendMessage(mq.BuildPage, mq.Frame, &dto.SwitchToPageCommand{Name: "ChaptersPage"}, false)
	// }
}
