package ui

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"abb_ia/internal/dto"
	"abb_ia/internal/mq"
	"abb_ia/internal/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type BuildPage struct {
	mq              *mq.Dispatcher
	grid            *tview.Grid
	infoPanel       *infoPanel
	buildSection    *tview.Grid
	copySection     *tview.Grid
	buildTable      *table
	copyTable       *table
	progressTable   *table
	progressSection *tview.Grid
}

func newBuildPage(dispatcher *mq.Dispatcher) *BuildPage {
	p := &BuildPage{}
	p.mq = dispatcher
	p.mq.RegisterListener(mq.BuildPage, p.dispatchMessage)

	p.grid = tview.NewGrid()
	p.grid.SetRows(7, -1, -1, 4)
	p.grid.SetColumns(0)

	// Ignore mouse events when the grid has no focus
	p.grid.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		if p.grid.HasFocus() {
			return action, event
		} else {
			return action, nil
		}
	})

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
	p.buildTable.setHeaders(" # ", "File name", "Format", "Duration", "Total Size", "Build progress")
	p.buildTable.setWeights(1, 2, 1, 1, 1, 5)
	p.buildTable.setAlign(tview.AlignRight, tview.AlignLeft, tview.AlignLeft, tview.AlignRight, tview.AlignRight, tview.AlignLeft)
	p.buildSection.AddItem(p.buildTable.t, 0, 0, 1, 1, 0, 0, true)
	p.grid.AddItem(p.buildSection, 1, 0, 1, 1, 0, 0, true)

	// copy section
	p.copySection = tview.NewGrid()
	p.copySection.SetColumns(-1)
	p.copySection.SetTitle(" Output directory copy progress: ")
	p.copySection.SetTitleAlign(tview.AlignLeft)
	p.copySection.SetBorder(true)

	p.copyTable = newTable()
	p.copyTable.setHeaders(" # ", "File name", "Format", "Duration", "Total Size", "Copy progress")
	p.copyTable.setWeights(1, 2, 1, 1, 1, 5)
	p.copyTable.setAlign(tview.AlignRight, tview.AlignLeft, tview.AlignLeft, tview.AlignRight, tview.AlignRight, tview.AlignLeft)
	p.copySection.AddItem(p.copyTable.t, 0, 0, 1, 1, 0, 0, true)
	p.grid.AddItem(p.copySection, 2, 0, 1, 1, 0, 0, true)

	// build progress section
	progressSection := tview.NewGrid()
	progressSection.SetColumns(-1)
	progressSection.SetBorder(true)
	progressSection.SetTitle(" Build progress: ")
	progressSection.SetTitleAlign(tview.AlignLeft)
	p.progressTable = newTable()
	p.progressTable.setWeights(1)
	p.progressTable.setAlign(tview.AlignLeft)
	p.progressTable.t.SetSelectable(false, false)
	progressSection.AddItem(p.progressTable.t, 0, 0, 1, 1, 0, 0, false)
	p.grid.AddItem(progressSection, 3, 0, 1, 1, 0, 0, false)
	p.progressSection = progressSection

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
	case *dto.BuildFileProgress:
		p.updateFileBuildProgress(dto)
	case *dto.BuildProgress:
		p.updateTotalBuildProgress(dto)
	case *dto.BuildComplete:
		p.buildComplete(dto)
	case *dto.UploadFileProgress:
		p.updateFileCopyProgress(dto)
	case *dto.CopyProgress:
		p.updateTotalCopyProgress(dto)
	case *dto.CopyComplete:
		p.copyComplete(dto)
	case *dto.UploadComplete:
		p.uploadComplete(dto)
	case *dto.ScanComplete:
		p.scanComplete(dto)
	case *dto.CleanupComplete:
		p.cleanupComplete(dto)
	default:
		m.UnsupportedTypeError(mq.BuildPage)
	}
}

func (p *BuildPage) displayBookInfo(ab *dto.Audiobook) {
	p.infoPanel.clear()
	// p.infoPanel.appendRow("", "")
	p.infoPanel.appendRow("Title:", ab.Title)
	p.infoPanel.appendRow("Author:", ab.Author)
	p.infoPanel.appendRow("Duration:", utils.SecondsToTime(ab.TotalDuration))
	p.infoPanel.appendRow("Size:", utils.BytesToHuman(ab.TotalSize))
	p.infoPanel.appendRow("Parts:", strconv.Itoa(len(ab.Parts)))

	p.buildTable.clear()
	p.buildTable.showHeader()
	for i, part := range ab.Parts {
		p.buildTable.appendRow(" "+strconv.Itoa(i+1)+" ", filepath.Base(part.M4BFile), part.Format, utils.SecondsToTime(part.Duration), utils.BytesToHuman(part.Size), "")
	}
	p.buildTable.ScrollToBeginning()

	if ab.Config.IsCopyToOutputDir() {
		p.copyTable.clear()
		p.copyTable.showHeader()
		for i, part := range ab.Parts {
			p.copyTable.appendRow(" "+strconv.Itoa(i+1)+" ", filepath.Base(part.M4BFile), part.Format, utils.SecondsToTime(part.Duration), utils.BytesToHuman(part.Size), "")
		}
		p.copyTable.ScrollToBeginning()
	}
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

func (p *BuildPage) updateFileBuildProgress(dp *dto.BuildFileProgress) {
	// update file progress
	col := 5
	w := p.buildTable.getColumnWidth(col) - 4
	progressText := fmt.Sprintf(" %3d%% ", dp.Percent)
	barWidth := int((float32((w - len(progressText))) * float32(dp.Percent) / 100))
	progressBar := strings.Repeat("━", barWidth) + strings.Repeat(" ", w-len(progressText)-barWidth)
	cell := p.buildTable.t.GetCell(dp.FileId+1, col)
	// cell.SetExpansion(0)
	// cell.SetMaxWidth(50)
	cell.Text = fmt.Sprintf("%s |%s|", progressText, progressBar)
	// p.buildTable.t.Select(dp.FileId+1, col)
	p.mq.SendMessage(mq.BuildPage, mq.TUI, &dto.DrawCommand{Primitive: nil}, false)
}

func (p *BuildPage) updateTotalBuildProgress(dp *dto.BuildProgress) {
	if p.progressTable.GetRowCount() == 0 {
		for i := 0; i < 2; i++ {
			p.progressTable.appendRow("")
		}
	}
	infoCell := p.progressTable.t.GetCell(0, 0)
	progressCell := p.progressTable.t.GetCell(1, 0)
	infoCell.Text = fmt.Sprintf("  [yellow]Time elapsed: [white]%10s | [yellow]Files: [white]%10s | [yellow]Speed: [white]%10s | [yellow]ETA: [white]%10s", dp.Elapsed, dp.Files, dp.Speed, dp.ETA)

	col := 0
	w := p.progressTable.getColumnWidth(col) - 4
	progressText := fmt.Sprintf(" %3d%% ", dp.Percent)
	barWidth := int((float32((w - len(progressText))) * float32(dp.Percent) / 100))
	progressBar := strings.Repeat("▒", barWidth) + strings.Repeat(" ", w-len(progressText)-barWidth)
	// progressCell.SetExpansion(0)
	// progressCell.SetMaxWidth(0)
	progressCell.Text = fmt.Sprintf("%s |%s|", progressText, progressBar)
	p.mq.SendMessage(mq.BuildPage, mq.TUI, &dto.DrawCommand{Primitive: nil}, false)
}

func (p *BuildPage) updateFileCopyProgress(dp *dto.UploadFileProgress) {
	// update file progress
	col := 5
	w := p.copyTable.getColumnWidth(col) - 3
	progressText := fmt.Sprintf(" %3d%% ", dp.Percent)
	barWidth := int((float32((w - len(progressText))) * float32(dp.Percent) / 100))
	progressBar := strings.Repeat("━", barWidth) + strings.Repeat(" ", w-len(progressText)-barWidth)
	cell := p.copyTable.t.GetCell(dp.FileId+1, col)
	// cell.SetExpansion(0)
	// cell.SetMaxWidth(50)
	cell.Text = fmt.Sprintf("%s |%s|", progressText, progressBar)
	// p.copyTable.t.Select(dp.FileId+1, col)
	p.mq.SendMessage(mq.BuildPage, mq.TUI, &dto.DrawCommand{Primitive: nil}, false)
}

func (p *BuildPage) updateTotalCopyProgress(dp *dto.CopyProgress) {
	if p.progressTable.GetRowCount() == 0 {
		for i := 0; i < 2; i++ {
			p.progressTable.appendRow("")
		}
	}
	p.progressSection.SetTitle(" Copy progress: ")
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
	p.mq.SendMessage(mq.BuildPage, mq.TUI, &dto.DrawCommand{Primitive: nil}, false)
}

/*
 * A chain of final operations -> ?Copy -> ?Upload -> ?Scan -> Cleanup - Done msg
 */
func (p *BuildPage) buildComplete(c *dto.BuildComplete) {
	// copy the book to Output directory if needed
	ab := c.Audiobook
	if ab.Config.IsCopyToOutputDir() {
		p.mq.SendMessage(mq.BuildPage, mq.CopyController, &dto.CopyCommand{Audiobook: ab}, true)
	} else {
		p.mq.SendMessage(mq.BuildPage, mq.BuildPage, &dto.CopyComplete{Audiobook: ab}, true)
	}
}

func (p *BuildPage) copyComplete(c *dto.CopyComplete) {
	// upload the book to Audiobookshelf server if needed
	ab := c.Audiobook
	if ab.Config.IsUploadToAudiobookshef() {
		p.mq.SendMessage(mq.BuildPage, mq.AudiobookshelfController, &dto.AudiobookshelfUploadCommand{Audiobook: ab}, true)
	} else {
		p.mq.SendMessage(mq.BuildPage, mq.BuildPage, &dto.AudiobookshelfScanCommand{Audiobook: ab}, true)
	}
}

func (p *BuildPage) uploadComplete(c *dto.UploadComplete) {
	// launch Audiobookshelf library scan if needed
	ab := c.Audiobook
	if ab.Config.IsScanAudiobookshef() {
		p.mq.SendMessage(mq.BuildPage, mq.AudiobookshelfController, &dto.AudiobookshelfScanCommand{Audiobook: c.Audiobook}, true)
	} else {
		p.mq.SendMessage(mq.BuildPage, mq.BuildPage, &dto.ScanComplete{Audiobook: ab}, true)
	}
}

func (p *BuildPage) scanComplete(c *dto.ScanComplete) {
	// clean up temporary directory
	p.mq.SendMessage(mq.BuildPage, mq.CleanupController, &dto.CleanupCommand{Audiobook: c.Audiobook}, true)
}

func (p *BuildPage) cleanupComplete(c *dto.CleanupComplete) {
	p.bookReadyMgs(c.Audiobook)
}

func (p *BuildPage) bookReadyMgs(ab *dto.Audiobook) {
	newMessageDialog(p.mq, "Build Complete", "Audiobook has been created", p.buildSection)
	//p.mq.SendMessage(mq.BuildPage, mq.Frame, &dto.SwitchToPageCommand{Name: "SearchPage"}, false)
}
