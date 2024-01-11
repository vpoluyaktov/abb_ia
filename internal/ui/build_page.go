package ui

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"abb_ia/internal/dto"
	"abb_ia/internal/mq"
	"abb_ia/internal/utils"

	"github.com/vpoluyaktov/tview"
)

type BuildPage struct {
	mq              *mq.Dispatcher
	mainGrid        *grid
	infoPanel       *infoPanel
	infoSection     *grid
	buildSection    *grid
	copySection     *grid
	uploadSection   *grid
	buildTable      *table
	copyTable       *table
	uploadTable     *table
	progressTable   *table
	progressSection *grid
	ab              *dto.Audiobook
}

func newBuildPage(dispatcher *mq.Dispatcher) *BuildPage {
	p := &BuildPage{}
	p.mq = dispatcher
	p.mq.RegisterListener(mq.BuildPage, p.dispatchMessage)

	p.mainGrid = newGrid()
	p.mainGrid.SetRows(7, -1, -1, -1, 4)
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

	// audiobook build section
	p.buildSection = newGrid()
	p.buildSection.SetColumns(-1)
	p.buildSection.SetTitle(" Audiobook build: ")
	p.buildSection.SetTitleAlign(tview.AlignLeft)
	p.buildSection.SetBorder(true)
	p.buildTable = newTable()
	p.buildTable.setHeaders(" # ", "File name", "Format", "Duration", "Size", "Build progress")
	p.buildTable.setWeights(1, 2, 1, 1, 1, 5)
	p.buildTable.setAlign(tview.AlignRight, tview.AlignLeft, tview.AlignLeft, tview.AlignRight, tview.AlignRight, tview.AlignLeft)
	p.buildSection.AddItem(p.buildTable.Table, 0, 0, 1, 1, 0, 0, true)
	p.mainGrid.AddItem(p.buildSection.Grid, 1, 0, 1, 1, 0, 0, true)

	// copy section
	p.copySection = newGrid()
	p.copySection.SetColumns(-1)
	p.copySection.SetTitle(" Copy to the output directory: ")
	p.copySection.SetTitleAlign(tview.AlignLeft)
	p.copySection.SetBorder(true)
	p.copyTable = newTable()
	p.copyTable.setHeaders(" # ", "File name", "Format", "Duration", "Size", "Copy Progress")
	p.copyTable.setWeights(1, 2, 1, 1, 1, 5)
	p.copyTable.setAlign(tview.AlignRight, tview.AlignLeft, tview.AlignLeft, tview.AlignRight, tview.AlignRight, tview.AlignLeft)
	p.copySection.AddItem(p.copyTable.Table, 0, 0, 1, 1, 0, 0, true)
	p.mainGrid.AddItem(p.copySection.Grid, 2, 0, 1, 1, 0, 0, true)

	// upload section
	p.uploadSection = newGrid()
	p.uploadSection.SetColumns(-1)
	p.uploadSection.SetTitle(" Upload to Audiobookshelf server: ")
	p.uploadSection.SetTitleAlign(tview.AlignLeft)
	p.uploadSection.SetBorder(true)
	p.uploadTable = newTable()
	p.uploadTable.setHeaders(" # ", "File name", "Format", "Duration", "Size", "Upload Progress")
	p.uploadTable.setWeights(1, 2, 1, 1, 1, 5)
	p.uploadTable.setAlign(tview.AlignRight, tview.AlignLeft, tview.AlignLeft, tview.AlignRight, tview.AlignRight, tview.AlignLeft)
	p.uploadSection.AddItem(p.uploadTable.Table, 0, 0, 1, 1, 0, 0, true)
	p.mainGrid.AddItem(p.uploadSection.Grid, 3, 0, 1, 1, 0, 0, true)

	// total progress section
	p.progressSection = newGrid()
	p.progressSection.SetColumns(-1)
	p.progressSection.SetBorder(true)
	p.progressSection.SetTitle(" Build progress: ")
	p.progressSection.SetTitleAlign(tview.AlignLeft)
	p.progressTable = newTable()
	p.progressTable.setWeights(1)
	p.progressTable.setAlign(tview.AlignLeft)
	p.progressTable.SetSelectable(false, false)
	p.progressSection.AddItem(p.progressTable.Table, 0, 0, 1, 1, 0, 0, false)
	p.mainGrid.AddItem(p.progressSection.Grid, 4, 0, 1, 1, 0, 0, false)

	p.mainGrid.SetNavigationOrder(
		p.infoPanel.Table,
		p.buildTable,
		p.copyTable,
		p.uploadTable,
		p.progressTable,
	)

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
	case *dto.FileBuildProgress:
		p.updateFileBuildProgress(dto)
	case *dto.TotalBuildProgress:
		p.updateTotalBuildProgress(dto)
	case *dto.BuildComplete:
		p.buildComplete(dto)
	case *dto.CopyFileProgress:
		p.updateFileCopyProgress(dto)
	case *dto.CopyProgress:
		p.updateTotalCopyProgress(dto)
	case *dto.UploadFileProgress:
		p.updateFileUploadProgress(dto)
	case *dto.UploadProgress:
		p.updateTotalUploadProgress(dto)
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

	p.ab = ab

	// dynamic grid layout generation
	p.mainGrid.Clear()
	p.mainGrid.SetColumns(0)
	p.mainGrid.SetRows(7, -1, 4)
	p.mainGrid.AddItem(p.infoSection.Grid, 0, 0, 1, 1, 0, 0, false)
	p.mainGrid.AddItem(p.buildSection.Grid, 1, 0, 1, 1, 0, 0, true)
	if ab.Config.IsCopyToOutputDir() && ab.Config.IsUploadToAudiobookshef() {
		p.mainGrid.SetRows(7, -1, -1, -1, 4)
		p.mainGrid.AddItem(p.copySection.Grid, 2, 0, 1, 1, 0, 0, true)
		p.mainGrid.AddItem(p.uploadSection.Grid, 3, 0, 1, 1, 0, 0, true)
		p.mainGrid.AddItem(p.progressSection.Grid, 4, 0, 1, 1, 0, 0, false)
	} else if ab.Config.IsCopyToOutputDir() {
		p.mainGrid.SetRows(7, -1, -1, 4)
		p.mainGrid.AddItem(p.copySection.Grid, 2, 0, 1, 1, 0, 0, true)
		p.mainGrid.AddItem(p.progressSection.Grid, 3, 0, 1, 1, 0, 0, false)
	} else if ab.Config.IsUploadToAudiobookshef() {
		p.mainGrid.SetRows(7, -1, -1, 4)
		p.mainGrid.AddItem(p.uploadSection.Grid, 2, 0, 1, 1, 0, 0, true)
		p.mainGrid.AddItem(p.progressSection.Grid, 3, 0, 1, 1, 0, 0, false)
	} else {
		p.mainGrid.AddItem(p.progressSection.Grid, 2, 0, 1, 1, 0, 0, false)
	}

	p.infoPanel.clear()
	// p.infoPanel.appendRow("", "")
	p.infoPanel.appendRow("Title:", ab.Title)
	p.infoPanel.appendRow("Author:", ab.Author)
	p.infoPanel.appendRow("Duration:", utils.SecondsToTime(ab.TotalDuration))
	p.infoPanel.appendRow("Size:", utils.BytesToHuman(ab.TotalSize))
	p.infoPanel.appendRow("Parts:", strconv.Itoa(len(ab.Parts)))

	p.buildTable.Clear()
	p.buildTable.showHeader()
	for i, part := range ab.Parts {
		p.buildTable.appendRow(" "+strconv.Itoa(i+1)+" ", filepath.Base(part.M4BFile), part.Format, utils.SecondsToTime(part.Duration), utils.BytesToHuman(part.Size), "")
	}
	p.buildTable.ScrollToBeginning()

	p.copyTable.Clear()
	p.copyTable.showHeader()
	for i, part := range ab.Parts {
		p.copyTable.appendRow(" "+strconv.Itoa(i+1)+" ", filepath.Base(part.M4BFile), part.Format, utils.SecondsToTime(part.Duration), utils.BytesToHuman(part.Size), "")
	}
	p.copyTable.ScrollToBeginning()

	p.uploadTable.Clear()
	p.uploadTable.showHeader()
	for i, part := range ab.Parts {
		p.uploadTable.appendRow(" "+strconv.Itoa(i+1)+" ", filepath.Base(part.M4BFile), part.Format, utils.SecondsToTime(part.Duration), utils.BytesToHuman(part.Size), "")
	}
	p.uploadTable.ScrollToBeginning()

	ui.SetFocus(p.buildTable.Table)
	ui.Draw()
}

func (p *BuildPage) stopConfirmation() {
	newYesNoDialog(p.mq, "Stop Confirmation", "Are you sure you want to stop the build?", p.buildSection.Grid, p.stopBuild, func() {})
}

func (p *BuildPage) stopBuild() {
	// Stop the build here
	p.mq.SendMessage(mq.BuildPage, mq.BuildController, &dto.StopCommand{Process: "Build", Reason: "User request"}, false)
	p.mq.SendMessage(mq.ChaptersPage, mq.CleanupController, &dto.CleanupCommand{Audiobook: p.ab}, true)
	p.switchToSearch()
}

func (p *BuildPage) updateFileBuildProgress(dp *dto.FileBuildProgress) {
	col := 5
	w := p.buildTable.GetColumnWidth(col) - 5
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
		progressBar := strings.Repeat(fileProgressChar, barWidth) + strings.Repeat(" ", fillerWidth)
		cell := p.buildTable.GetCell(dp.FileId+1, col)
		cell.SetExpansion(p.buildTable.colWeight[col])
		cell.Text = fmt.Sprintf("%s |%s|", progressText, progressBar)
		ui.Draw()
	}
}

func (p *BuildPage) updateTotalBuildProgress(dp *dto.TotalBuildProgress) {
	if p.progressTable.GetRowCount() == 0 {
		for i := 0; i < 2; i++ {
			p.progressTable.appendRow("")
		}
	}
	p.progressSection.SetTitle(" Build progress: ")
	infoCell := p.progressTable.GetCell(0, 0)
	progressCell := p.progressTable.GetCell(1, 0)
	infoCell.Text = fmt.Sprintf("  [yellow]Time elapsed: [white]%10s | [yellow]Files: [white]%10s | [yellow]Speed: [white]%10s | [yellow]ETA: [white]%10s", dp.Elapsed, dp.Files, dp.Speed, dp.ETA)

	col := 0
	w := p.progressTable.GetColumnWidth(col) - 4
	if w > 0 {
		progressText := fmt.Sprintf(" %3d%% ", dp.Percent)
		barWidth := int((float32((w - len(progressText))) * float32(dp.Percent) / 100))
		progressBar := strings.Repeat(totalProgressChar, barWidth) + strings.Repeat(" ", w-len(progressText)-barWidth)
		progressCell.Text = fmt.Sprintf("%s |%s|", progressText, progressBar)
		ui.Draw()
	}
}

func (p *BuildPage) updateFileCopyProgress(dp *dto.CopyFileProgress) {
	// update file progress
	col := 5
	w := p.copyTable.GetColumnWidth(col) - 5
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
		progressBar := strings.Repeat(fileProgressChar, barWidth) + strings.Repeat(" ", fillerWidth)
		cell := p.copyTable.GetCell(dp.FileId+1, col)
		cell.SetExpansion(p.copyTable.colWeight[col])
		cell.Text = fmt.Sprintf("%s |%s|", progressText, progressBar)
		ui.Draw()
	}
}

func (p *BuildPage) updateTotalCopyProgress(dp *dto.CopyProgress) {
	if p.progressTable.GetRowCount() == 0 {
		for i := 0; i < 2; i++ {
			p.progressTable.appendRow("")
		}
	}
	p.progressSection.SetTitle(" Copy progress: ")
	infoCell := p.progressTable.GetCell(0, 0)
	progressCell := p.progressTable.GetCell(1, 0)
	infoCell.Text = fmt.Sprintf("  [yellow]Time elapsed: [white]%10s | [yellow]Files: [white]%10s | [yellow]Speed: [white]%12s | [yellow]ETA: [white]%10s", dp.Elapsed, dp.Files, dp.Speed, dp.ETA)

	col := 0
	w := p.progressTable.GetColumnWidth(col) - 5
	if w > 0 {
		progressText := fmt.Sprintf(" %3d%% ", dp.Percent)
		barWidth := int((float32((w - len(progressText))) * float32(dp.Percent) / 100))
		progressBar := strings.Repeat(totalProgressChar, barWidth) + strings.Repeat(" ", w-len(progressText)-barWidth)
		// progressCell.SetExpansion(0)
		// progressCell.SetMaxWidth(0)
		progressCell.Text = fmt.Sprintf("%s |%s|", progressText, progressBar)
		ui.Draw()
	}
}

func (p *BuildPage) updateFileUploadProgress(dp *dto.UploadFileProgress) {
	col := 5
	w := p.uploadTable.GetColumnWidth(col) - 5
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
		progressBar := strings.Repeat(fileProgressChar, barWidth) + strings.Repeat(" ", fillerWidth)
		cell := p.uploadTable.GetCell(dp.FileId+1, col)
		cell.SetExpansion(p.uploadTable.colWeight[col])
		cell.Text = fmt.Sprintf("%s |%s|", progressText, progressBar)
		ui.Draw()
	}
}

func (p *BuildPage) updateTotalUploadProgress(dp *dto.UploadProgress) {
	if p.progressTable.GetRowCount() == 0 {
		for i := 0; i < 2; i++ {
			p.progressTable.appendRow("")
		}
	}
	p.progressSection.SetTitle(" Upload progress: ")
	infoCell := p.progressTable.GetCell(0, 0)
	progressCell := p.progressTable.GetCell(1, 0)
	infoCell.Text = fmt.Sprintf("  [yellow]Time elapsed: [white]%10s | [yellow]Files: [white]%10s | [yellow]Speed: [white]%12s | [yellow]ETA: [white]%10s", dp.Elapsed, dp.Files, dp.Speed, dp.ETA)

	col := 0
	w := p.progressTable.GetColumnWidth(col) - 5
	if w > 0 {
		progressText := fmt.Sprintf(" %3d%% ", dp.Percent)
		barWidth := int((float32((w - len(progressText))) * float32(dp.Percent) / 100))
		progressBar := strings.Repeat(totalProgressChar, barWidth) + strings.Repeat(" ", w-len(progressText)-barWidth)
		// progressCell.SetExpansion(0)
		// progressCell.SetMaxWidth(0)
		progressCell.Text = fmt.Sprintf("%s |%s|", progressText, progressBar)
		ui.Draw()
	}
}

/*
 * A chain of final operations: ?Copy -> ?Upload -> ?Scan -> Cleanup - Done msg
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
		p.mq.SendMessage(mq.BuildPage, mq.UploadController, &dto.AbsUploadCommand{Audiobook: ab}, true)
	} else {
		p.mq.SendMessage(mq.BuildPage, mq.BuildPage, &dto.UploadComplete{Audiobook: ab}, true)
	}
}

func (p *BuildPage) uploadComplete(c *dto.UploadComplete) {
	// launch Audiobookshelf library scan if needed
	ab := c.Audiobook
	if ab.Config.IsScanAudiobookshef() {
		p.mq.SendMessage(mq.BuildPage, mq.UploadController, &dto.AbsScanCommand{Audiobook: c.Audiobook}, true)
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
	newMessageDialog(p.mq, "Build Complete", "Audiobook has been created", p.buildSection.Grid, p.switchToSearch)
}

func (p *BuildPage) switchToSearch() {
	p.mq.SendMessage(mq.BuildPage, mq.Frame, &dto.SwitchToPageCommand{Name: "SearchPage"}, false)
}
