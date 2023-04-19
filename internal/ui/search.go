package ui

import (
	"strconv"

	"github.com/rivo/tview"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/dto"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/logger"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/mq"
)

type searchPanel struct {
	dispatcher     *mq.Dispatcher
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

func newSearchPanel(dispatcher *mq.Dispatcher) *searchPanel {
	p := &searchPanel{}
	p.dispatcher = dispatcher
	p.dispatcher.RegisterListener(mq.SearchPanel, p.dispatchMessage)

	p.grid = tview.NewGrid()
	p.grid.SetRows(5, -1, -1)
	p.grid.SetColumns(0)

	// search section
	p.searchSection = tview.NewGrid()
	p.searchSection.SetColumns(-3, -1)
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
	p.searchButton = f.AddButton("Create Audiobook", p.createBook)
	p.clearButton = f.AddButton("Settings", p.clearEverything)
	p.searchSection.AddItem(f.f, 0, 1, 1, 1, 0, 0, false)

	p.grid.AddItem(p.searchSection, 0, 0, 1, 1, 0, 0, true)

	// result section
	p.resultSection = tview.NewGrid()
	p.resultSection.SetColumns(-1)
	p.resultSection.SetTitle(" Search result ")
	p.resultSection.SetTitleAlign(tview.AlignLeft)
	p.resultSection.SetBorder(true)

	p.resultTable = newTable()
	p.resultTable.setHeaders("Title", "Files", "Duration (HH:MM:SS)", "Total Size")
	p.resultTable.setWidths(6, 2, 1, 1)
	p.resultTable.setAlign(tview.AlignLeft, tview.AlignRight, tview.AlignRight, tview.AlignRight)
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

	p.sendMessage(mq.SearchPanel, mq.TUI, dto.SetFocusCommandType, &dto.SetFocusCommand{Primitive: p.searchSection}, true)

	return p
}

func (p *searchPanel) checkMQ() {
	m := p.dispatcher.GetMessage(mq.SearchPanel)
	if m != nil {
		p.dispatchMessage(m)
	}
}

func (p *searchPanel) sendMessage(from string, to string, dtoType string, dto dto.Dto, async bool) {
	m := &mq.Message{}
	m.From = from
	m.To = to
	m.Type = dtoType
	m.Dto = dto
	m.Async = async
	p.dispatcher.SendMessage(m)
}

func (p *searchPanel) dispatchMessage(m *mq.Message) {
	switch t := m.Type; t {
	case dto.IAItemType:
		if r, ok := m.Dto.(*dto.IAItem); ok {
			go p.updateResult(r)
		} else {
			m.DtoCastError()
		}

	default:
		m.UnsupportedTypeError()
	}
}

func (p *searchPanel) runSearch() {
	p.clearSearchResults()
	p.resultTable.showHeader()
	// Disable Search Button here
	p.sendMessage(mq.SearchPanel, mq.SearchController, dto.SearchCommandType, &dto.SearchCommand{SearchCondition: p.searchCriteria}, true)
	p.sendMessage(mq.SearchPanel, mq.TUI, dto.SetFocusCommandType, &dto.SetFocusCommand{Primitive: p.resultTable.t}, true)
}

func (p *searchPanel) clearSearchResults() {
	p.searchResult = make([]*dto.IAItem, 0)
	p.resultTable.clear()
	p.descriptionView.SetText("")
	p.filesTable.clear()
}

func (p *searchPanel) clearEverything() {
	p.inputField.SetText("")
	p.clearSearchResults()
}

func (p *searchPanel) updateResult(i *dto.IAItem) {
	logger.Debug(mq.SearchPanel + ": Got AI Item: " + i.Title)
	p.searchResult = append(p.searchResult, i)
	p.resultTable.appendRow(i.Title, strconv.Itoa(i.FilesCount), i.TotalLengthH, i.TotalSizeH)
	p.resultTable.t.ScrollToBeginning()
	p.sendMessage(mq.SearchPanel, mq.TUI, dto.DrawCommandType, &dto.DrawCommand{Primitive: p.resultTable.t}, true)
	p.updateDetails(1, 0)

}

func (p *searchPanel) updateDetails(row int, col int) {
	if row > 0 && len(p.searchResult) > 0 && len(p.searchResult) >= row {
		d := p.searchResult[row-1].Description
		p.descriptionView.SetText(d)
		p.descriptionView.ScrollToBeginning()
		p.sendMessage(mq.SearchPanel, mq.TUI, dto.DrawCommandType, &dto.DrawCommand{Primitive: p.descriptionView}, true)

		p.filesTable.clear()
		p.filesTable.showHeader()
		files := p.searchResult[row-1].Files
		for _, f := range files {
			p.filesTable.appendRow(f.Name, f.Format, f.LengthH, f.SizeH)
		}
		p.filesTable.t.ScrollToBeginning()
		// p.sendMessage(mq.SearchPanel, mq.TUI, dto.DrawCommandType, &dto.DrawCommand{Primitive: p.filesTable.t}, true)
		p.sendMessage(mq.SearchPanel, mq.TUI, dto.DrawCommandType, &dto.DrawCommand{Primitive: nil}, true)
	}

}

func (p *searchPanel) createBook() {
	m := newDialogWindow(p.dispatcher, 25, 60)
	m.Show()
}
