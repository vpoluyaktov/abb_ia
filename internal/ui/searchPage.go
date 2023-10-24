package ui

import (
	"fmt"
	"strconv"

	"github.com/rivo/tview"
	"github.com/vpoluyaktov/abb_ia/internal/config"
	"github.com/vpoluyaktov/abb_ia/internal/dto"
	"github.com/vpoluyaktov/abb_ia/internal/logger"
	"github.com/vpoluyaktov/abb_ia/internal/mq"
	"github.com/vpoluyaktov/abb_ia/internal/utils"
)

type SearchPage struct {
	mq             *mq.Dispatcher
	grid           *tview.Grid
	searchCriteria string
	searchResult   []*dto.IAItem

	searchSection *tview.Grid
	inputField    *tview.InputField
	searchButton  *tview.Button
	clearButton   *tview.Button

	resultSection *tview.Grid
	resultTable   *table

	detailsSection  *tview.Grid
	descriptionView *tview.TextView
	filesTable      *table
}

func newSearchPage(dispatcher *mq.Dispatcher) *SearchPage {
	p := &SearchPage{}
	p.mq = dispatcher
	p.mq.RegisterListener(mq.SearchPage, p.dispatchMessage)

	p.searchCriteria = config.SearchCondition()

	p.grid = tview.NewGrid()
	p.grid.SetRows(5, -1, -1)
	p.grid.SetColumns(0)

	// search section
	p.searchSection = tview.NewGrid()
	p.searchSection.SetColumns(-2, -1)
	p.searchSection.SetBorder(true)
	p.searchSection.SetTitle(" Internet Archive Search ")
	p.searchSection.SetTitleAlign(tview.AlignLeft)
	f := newForm()
	f.SetHorizontal(true)
	p.inputField = f.AddInputField("Search criteria", config.SearchCondition(), 40, nil, func(t string) { p.searchCriteria = t })
	p.searchButton = f.AddButton("Search", p.runSearch)
	p.clearButton = f.AddButton("Clear", p.clearEverything)
	p.searchSection.AddItem(f.f, 0, 0, 1, 1, 0, 0, true)
	f = newForm()
	f.SetHorizontal(false)
	f.f.SetButtonsAlign(tview.AlignRight)
	p.searchButton = f.AddButton("Create Audiobook", p.createBook)
	p.clearButton = f.AddButton("Settings", p.updateConfig)
	p.searchSection.AddItem(f.f, 0, 1, 1, 1, 0, 0, false)

	p.grid.AddItem(p.searchSection, 0, 0, 1, 1, 0, 0, true)

	// result section
	p.resultSection = tview.NewGrid()
	p.resultSection.SetColumns(-1)
	p.resultSection.SetTitle(" Search result: ")
	p.resultSection.SetTitleAlign(tview.AlignLeft)
	p.resultSection.SetBorder(true)

	p.resultTable = newTable()
	p.resultTable.setHeaders("Author", "Title", "Files", "Duration (hh:mm:ss)", "Total Size")
	p.resultTable.setWeights(3, 6, 2, 1, 1)
	p.resultTable.setAlign(tview.AlignLeft, tview.AlignLeft, tview.AlignRight, tview.AlignRight, tview.AlignRight)
	p.resultTable.t.SetSelectionChangedFunc(p.updateDetails)
	p.resultSection.AddItem(p.resultTable.t, 0, 0, 1, 1, 0, 0, true)
	p.grid.AddItem(p.resultSection, 1, 0, 1, 1, 0, 0, true)

	// details section
	p.detailsSection = tview.NewGrid()
	p.detailsSection.SetRows(-1)
	p.detailsSection.SetColumns(-1, 1, -1)

	p.descriptionView = tview.NewTextView()
	p.descriptionView.SetWrap(true)
	p.descriptionView.SetWordWrap(true)
	p.descriptionView.SetBorder(true)
	p.descriptionView.SetTitle(" Description: ")
	p.descriptionView.SetTitleAlign(tview.AlignLeft)
	p.detailsSection.AddItem(p.descriptionView, 0, 0, 1, 1, 0, 0, true)

	p.filesTable = newTable()
	p.filesTable.t.SetBorder(true)
	p.filesTable.t.SetTitle(" Files: ")
	p.filesTable.t.SetTitleAlign(tview.AlignLeft)
	p.filesTable.setHeaders("File name", "Format", "Duration", "Size")
	p.filesTable.setWeights(3, 1, 1, 1)
	p.filesTable.setAlign(tview.AlignLeft, tview.AlignRight, tview.AlignRight, tview.AlignRight)
	p.detailsSection.AddItem(p.filesTable.t, 0, 2, 1, 1, 0, 0, true)

	p.grid.AddItem(p.detailsSection, 2, 0, 1, 1, 0, 0, true)

	p.mq.SendMessage(mq.SearchPage, mq.Frame, &dto.SwitchToPageCommand{Name: "SearchPage"}, false)
	p.mq.SendMessage(mq.SearchPage, mq.TUI, &dto.SetFocusCommand{Primitive: p.searchSection}, true)

	return p
}

func (p *SearchPage) checkMQ() {
	m := p.mq.GetMessage(mq.SearchPage)
	if m != nil {
		p.dispatchMessage(m)
	}
}

func (p *SearchPage) dispatchMessage(m *mq.Message) {
	switch dto := m.Dto.(type) {
	case *dto.IAItem:
		go p.updateResult(dto)
	case *dto.SearchProgress:
		p.updateTitle(dto)
	default:
		m.UnsupportedTypeError(mq.SearchPage)
	}
}

func (p *SearchPage) runSearch() {
	p.clearSearchResults()
	p.resultTable.showHeader()
	// Disable Search Button here
	p.mq.SendMessage(mq.SearchPage, mq.SearchController, &dto.SearchCommand{SearchCondition: p.searchCriteria}, false)
	p.mq.SendMessage(mq.SearchPage, mq.TUI, &dto.SetFocusCommand{Primitive: p.resultTable.t}, true)
}

func (p *SearchPage) clearSearchResults() {
	p.searchResult = make([]*dto.IAItem, 0)
	p.resultSection.SetTitle(" Search result: ")
	p.resultTable.clear()
	p.descriptionView.SetText("")
	p.filesTable.clear()
}

func (p *SearchPage) clearEverything() {
	p.inputField.SetText("")
	p.clearSearchResults()
}

func (p *SearchPage) updateResult(i *dto.IAItem) {
	logger.Debug(mq.SearchPage + ": Got AI Item: " + i.Title)
	p.searchResult = append(p.searchResult, i)
	p.resultTable.appendRow(i.Creator, i.Title, strconv.Itoa(len(i.AudioFiles)), utils.SecondsToTime(i.TotalLength), utils.BytesToHuman(i.TotalSize))
	p.resultTable.ScrollToBeginning()
	// p.mq.SendMessage(mq.SearchPage, mq.TUI, &dto.DrawCommand{Primitive: p.resultTable.t}, true) // single primitive refresh is not supported by tview (but supported by cview)
	p.updateDetails(1, 0)
	p.mq.SendMessage(mq.SearchPage, mq.TUI, &dto.DrawCommand{Primitive: nil}, true)
}

func (p *SearchPage) updateTitle(sp *dto.SearchProgress) {
	p.resultSection.SetTitle(fmt.Sprintf(" Search result (first %d items from %d total): ", sp.ItemsFetched, sp.ItemsTotal))
}

func (p *SearchPage) updateDetails(row int, col int) {
	if row > 0 && len(p.searchResult) > 0 && len(p.searchResult) >= row {
		d := p.searchResult[row-1].Description
		p.descriptionView.SetText(d)
		p.descriptionView.ScrollToBeginning()
		// p.mq.SendMessage(mq.SearchPage, mq.TUI, &dto.DrawCommand{Primitive: p.descriptionView}, true) // single primitive refresh is not supported by tview (but supported by cview)

		p.filesTable.clear()
		p.filesTable.showHeader()
		files := p.searchResult[row-1].AudioFiles
		for _, f := range files {
			p.filesTable.appendRow(f.Name, f.Format, utils.SecondsToTime(f.Length), utils.BytesToHuman(f.Size))
		}
		p.filesTable.ScrollToBeginning()
		// p.mq.SendMessage(mq.SearchPage, mq.TUI, &dto.DrawCommand{Primitive: p.filesTable.t}, true) // single primitive refresh is not supported by tview (but supported by cview)
	}
}

func (p *SearchPage) createBook() {
	// get selectet row from the results table
	row, _ := p.resultTable.t.GetSelection()
	if row <= 0 || len(p.searchResult) <= 0 || len(p.searchResult) < row {
		newMessageDialog(p.mq, "Error", "Please perform a search first", p.searchSection)
	} else {
		item := p.searchResult[row-1]
		d := newDialogWindow(p.mq, 12, 80, p.resultSection)
		f := newForm()
		f.SetTitle("Create Audiobook")
		// author := f.AddInputField("Book Author", item.Creator, 60, nil, nil)
		// title := f.AddInputField("Book Title", item.Title, 60, nil, nil)
		f.AddButton("Create Audiobook", func() {
			ab := &dto.Audiobook{}
			ab.IAItem = item
			// ab.Title = title.GetText()
			// ab.Author = author.GetText()
			p.startDownload(ab)
			d.Close()
		})
		f.AddButton("Cancel", func() {
			d.Close()
		})
		d.setForm(f.f)
		d.Show()
	}
}

func (p *SearchPage) startDownload(ab *dto.Audiobook) {
	p.mq.SendMessage(mq.SearchPage, mq.DownloadController, &dto.DownloadCommand{Audiobook: ab}, true)
	p.mq.SendMessage(mq.SearchPage, mq.Frame, &dto.SwitchToPageCommand{Name: "DownloadPage"}, false)
}

func (p *SearchPage) updateConfig() {
	p.mq.SendMessage(mq.SearchPage, mq.ConfigPage, &dto.DisplayConfigCommand{Config: config.GetCopy()}, true)
	p.mq.SendMessage(mq.SearchPage, mq.Frame, &dto.SwitchToPageCommand{Name: "ConfigPage"}, false)
}
