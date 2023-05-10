package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rivo/tview"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/dto"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/logger"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/mq"
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
	p.inputField = f.AddInputField("Search criteria", "", 40, nil, func(t string) { p.searchCriteria = t })
	p.searchButton = f.AddButton("Search", p.runSearch)
	p.clearButton = f.AddButton("Clear", p.clearEverything)
	p.searchSection.AddItem(f.f, 0, 0, 1, 1, 0, 0, true)
	f = newForm()
	f.SetHorizontal(false)
	f.f.SetButtonsAlign(tview.AlignRight)
	p.searchButton = f.AddButton("Create Audiobook", p.createBook)
	p.clearButton = f.AddButton("Settings", p.clearEverything)
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
	p.resultTable.setWidths(3, 6, 2, 1, 1)
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
	p.filesTable.setWidths(3, 1, 1, 1)
	p.filesTable.setAlign(tview.AlignLeft, tview.AlignRight, tview.AlignRight, tview.AlignRight)
	p.detailsSection.AddItem(p.filesTable.t, 0, 2, 1, 1, 0, 0, true)

	p.grid.AddItem(p.detailsSection, 2, 0, 1, 1, 0, 0, true)

	p.mq.SendMessage(mq.SearchPage, mq.TUI, dto.SetFocusCommandType, &dto.SetFocusCommand{Primitive: p.searchSection}, true)

	return p
}

func (p *SearchPage) checkMQ() {
	m := p.mq.GetMessage(mq.SearchPage)
	if m != nil {
		p.dispatchMessage(m)
	}
}

func (p *SearchPage) dispatchMessage(m *mq.Message) {
	switch t := m.Type; t {
	case dto.IAItemType:
		if r, ok := m.Dto.(*dto.IAItem); ok {
			go p.updateResult(r)
		} else {
			m.DtoCastError(mq.SearchPage)
		}
	case dto.SearchProgressType:
		if sp, ok := m.Dto.(*dto.SearchProgress); ok {
			p.updateTitle(sp)
		} else {
			m.DtoCastError(mq.SearchPage)
		}

	default:
		m.UnsupportedTypeError(mq.SearchPage)
	}
}

func (p *SearchPage) runSearch() {
	p.clearSearchResults()
	p.resultTable.showHeader()
	// Disable Search Button here
	p.mq.SendMessage(mq.SearchPage, mq.SearchController, dto.SearchCommandType, &dto.SearchCommand{SearchCondition: p.searchCriteria}, false)
	p.mq.SendMessage(mq.SearchPage, mq.TUI, dto.SetFocusCommandType, &dto.SetFocusCommand{Primitive: p.resultTable.t}, true)
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
	p.resultTable.appendRow(i.Creator, i.Title, strconv.Itoa(i.FilesCount), i.TotalLengthH, i.TotalSizeH)
	p.resultTable.t.ScrollToBeginning()
	// p.mq.SendMessage(mq.SearchPage, mq.TUI, dto.DrawCommandType, &dto.DrawCommand{Primitive: p.resultTable.t}, true) // single primitive refresh is not supported by tview (but supported by cview)
	p.updateDetails(1, 0)
	p.mq.SendMessage(mq.SearchPage, mq.TUI, dto.DrawCommandType, &dto.DrawCommand{Primitive: nil}, true)
}

func (p *SearchPage) updateTitle(sp *dto.SearchProgress) {
	p.resultSection.SetTitle(fmt.Sprintf(" Search result (first %d items from %d total): ", sp.ItemsFetched, sp.ItemsTotal))
}

func (p *SearchPage) updateDetails(row int, col int) {
	if row > 0 && len(p.searchResult) > 0 && len(p.searchResult) >= row {
		d := p.searchResult[row-1].Description
		p.descriptionView.SetText(d)
		p.descriptionView.ScrollToBeginning()
		// p.mq.SendMessage(mq.SearchPage, mq.TUI, dto.DrawCommandType, &dto.DrawCommand{Primitive: p.descriptionView}, true) // single primitive refresh is not supported by tview (but supported by cview)

		p.filesTable.clear()
		p.filesTable.showHeader()
		files := p.searchResult[row-1].Files
		for _, f := range files {
			p.filesTable.appendRow(strings.TrimPrefix(f.Name, "/"), f.Format, f.LengthH, f.SizeH)
		}
		p.filesTable.t.ScrollToBeginning()
		// p.mq.SendMessage(mq.SearchPage, mq.TUI, dto.DrawCommandType, &dto.DrawCommand{Primitive: p.filesTable.t}, true) // single primitive refresh is not supported by tview (but supported by cview)
	}

}

func (p *SearchPage) createBook() {
	// get selectet row from the results table
	row, _ := p.resultTable.t.GetSelection()
	if row > 0 && len(p.searchResult) > 0 && len(p.searchResult) >= row {
		item := p.searchResult[row-1]

		d := newDialogWindow(p.mq, 12, 80)
		f := newForm()
		f.SetTitle("Create Audiobook")
		f.AddInputField("Book Author", item.Creator, 60, nil, nil)
		f.AddInputField("Book Title", item.Title, 60, nil, nil)
		f.AddButton("Create Audiobook", func() {
			d.Close()
			p.launchDownload()
		})
		f.AddButton("Cancel", func() {
			d.Close()
		})
		d.setForm(f.f)
		d.Show()
	} else {
		newMessageDialog(p.mq, "Error", "Please perform a search first")
	}
}

func (p *SearchPage) launchDownload() {
	// d.Close()
	// p.mq.SendMessage(mq.SearchPage, mq.DownloadPageController, dto.SearchCommandType, &dto.SearchCommand{SearchCondition: p.searchCriteria}, true)
	p.mq.SendMessage(mq.SearchPage, mq.Frame, dto.SwitchToPageCommandType, &dto.SwitchToPageCommand{Name: "DownloadPage"}, false)
}
