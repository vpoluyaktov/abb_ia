package ui

import (
	"strconv"

	"code.rocketnine.space/tslocum/cview"

	"github.com/vpoluyaktov/audiobook_creator_IA/internal/dto"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/logger"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/mq"
)

const (
	uiComponentName = "SearchPanel"
	controllerName  = "SearchController"
)

type searchPanel struct {
	grid           *cview.Grid
	dispatcher     *mq.Dispatcher
	searchCriteria string
	searchResult   []*dto.IAItem

	searchSection *cview.Flex
	inputField    *cview.InputField
	searchButton  *cview.Button
	clearButton   *cview.Button

	resultSection *cview.Grid
	resultTable   *table

	detailsSection  *cview.Grid
	descriptionView *cview.TextView
	filesTable      *table
}

func newSearchPanel(dispatcher *mq.Dispatcher) *searchPanel {
	p := &searchPanel{}
	p.dispatcher = dispatcher
	p.dispatcher.RegisterListener(uiComponentName, p.dispatchMessage)

	p.grid = cview.NewGrid()
	p.grid.SetRows(5, -1, -1)
	p.grid.SetColumns(0)

	// search section
	p.searchSection = cview.NewFlex()
	p.searchSection.SetDirection(cview.FlexRow)
	p.searchSection.SetBorder(true)
	p.searchSection.SetTitle(" Internet Archive Search ")
	p.searchSection.SetTitleAlign(cview.AlignLeft)
	f := newForm()
	f.SetHorizontal(true)
	p.inputField = f.AddInputField("Search criteria", "", 40, nil, func(t string) { p.searchCriteria = t })
	p.searchButton = f.AddButton("Search", p.runSearch)
	p.clearButton = f.AddButton("Clear", p.clearEverything)
	p.searchSection.AddItem(f.f, 0, 1, true)
	p.grid.AddItem(p.searchSection, 0, 0, 1, 1, 0, 0, true)

	// result section
	p.resultSection = cview.NewGrid()
	p.resultSection.SetColumns(-1)
	p.resultSection.SetTitle(" Search result ")
	p.resultSection.SetTitleAlign(cview.AlignLeft)
	p.resultSection.SetBorder(true)

	p.resultTable = newTable()
	p.resultTable.setHeaders("Title", "Files", "Duration (HH:MM:SS)", "Total Size")
	p.resultTable.setWidths(6, 2, 1, 1)
	p.resultTable.setAlign(cview.AlignLeft, cview.AlignRight, cview.AlignRight, cview.AlignRight)
	p.resultTable.t.SetSelectionChangedFunc(p.updateDetails)
	p.resultSection.AddItem(p.resultTable.t, 0, 0, 1, 1, 0, 0, true)
	p.grid.AddItem(p.resultSection, 1, 0, 1, 1, 0, 0, true)

	// details section
	p.detailsSection = cview.NewGrid()
	p.detailsSection.SetRows(-1)
	p.detailsSection.SetColumns(-1, 1, -1)

	p.descriptionView = cview.NewTextView()
	p.descriptionView.SetWrap(true)
	p.descriptionView.SetWordWrap(true)
	p.descriptionView.SetBorder(true)
	p.descriptionView.SetTitle(" Description: ")
	p.descriptionView.SetTitleAlign(cview.AlignLeft)
	p.detailsSection.AddItem(p.descriptionView, 0, 0, 1, 1, 0, 0, true)

	p.filesTable = newTable()
	p.filesTable.t.SetBorder(true)
	p.filesTable.t.SetTitle(" Files: ")
	p.filesTable.t.SetTitleAlign(cview.AlignLeft)
	p.filesTable.setHeaders("File name", "Format", "Duration", "Size")
	p.filesTable.setWidths(3, 1, 1, 1)
	p.filesTable.setAlign(cview.AlignLeft, cview.AlignRight, cview.AlignRight, cview.AlignRight)
	p.detailsSection.AddItem(p.filesTable.t, 0, 2, 1, 1, 0, 0, true)

	p.grid.AddItem(p.detailsSection, 2, 0, 1, 1, 0, 0, true)

	p.sendMessage(uiComponentName, "TUI", dto.SetFocusCommandType, dto.SetFocusCommand{Primitive: p.searchSection}, true)

	return p
}

func (p *searchPanel) checkMQ() {
	m := p.dispatcher.GetMessage(uiComponentName)
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
	p.sendMessage(uiComponentName, controllerName, dto.SearchCommandType, dto.SearchCommand{SearchCondition: p.searchCriteria}, true)
	p.sendMessage(uiComponentName, "TUI", dto.SetFocusCommandType, dto.SetFocusCommand{Primitive: p.resultSection}, true)
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
	logger.Debug(uiComponentName + ": Got AI Item: " + i.Title)
	p.searchResult = append(p.searchResult, i)
	p.resultTable.appendRow(i.Title, strconv.Itoa(i.FilesCount), i.TotalLengthH, i.TotalSizeH)
	p.updateDetails(1, 0)
	p.sendMessage(uiComponentName, "TUI", dto.GeneralCommandType, dto.GeneralCommand{Command: "RedrawUI"}, true)
}

func (p *searchPanel) updateDetails(row int, col int) {
	if row > 0 && len(p.searchResult) > 0 && len(p.searchResult) >= row {
		d := p.searchResult[row-1].Description
		p.descriptionView.SetText(d)
		p.descriptionView.ScrollToBeginning()

		p.filesTable.clear()
		p.filesTable.showHeader()
		files := p.searchResult[row-1].Files
		for _, f := range files {
			p.filesTable.appendRow(f.Name, f.Format, f.LengthH, f.SizeH)
		}
		p.filesTable.t.ScrollToBeginning()
	}

}
