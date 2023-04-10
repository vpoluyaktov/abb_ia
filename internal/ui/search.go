package ui

import (
	"strconv"

	"code.rocketnine.space/tslocum/cview"
	"github.com/gdamore/tcell/v2"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/dto"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/logger"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/mq"
)

const (
	uiComponentName = "SearchPanel"
	controllerName  = "SearchController"
)

type searchPanel struct {
	grid              *cview.Grid
	dispatcher        *mq.Dispatcher
	searchCriteria    string
	searchResult      []*dto.IAItem
	searchResultSection *cview.Grid
	searchResultTable *cview.Table
	descriptionView   *cview.TextView
}

func newSearchPanel(dispatcher *mq.Dispatcher) *searchPanel {
	p := &searchPanel{}
	p.dispatcher = dispatcher
	p.dispatcher.RegisterListener(uiComponentName, p.dispatchMessage)

	p.grid = cview.NewGrid()
	p.grid.SetRows(5, -1, -1)
	p.grid.SetColumns(0)

	searchSection := cview.NewFlex()
	searchSection.SetDirection(cview.FlexRow)
	searchSection.SetBorder(true)
	searchSection.SetTitle(" Internet Archive Search ")
	searchSection.SetTitleAlign(cview.AlignLeft)
	form := cview.NewForm()
	form.SetHorizontal(true)
	form.AddInputField("Search criteria", "", 40, nil, func(t string) { p.searchCriteria = t })
	form.AddButton("Search", p.runSearch)
	form.AddButton("Clear", p.clearSearchResults)
	searchSection.AddItem(form, 0, 1, true)
	p.grid.AddItem(searchSection, 0, 0, 1, 1, 0, 0, true)

	// p.searchResultSection := cview.NewFlex()
	// p.searchResultSection.SetDirection(cview.FlexRow)
	p.searchResultSection = cview.NewGrid()
	p.searchResultSection.SetColumns(-1)
	p.searchResultSection.SetTitle(" Search result ")
	p.searchResultSection.SetTitleAlign(cview.AlignLeft)
	p.searchResultSection.SetBorder(true)
	p.searchResultTable = cview.NewTable()
	p.searchResultTable.SetSelectable(true, false)
	// p.searchResultTable.SetScrollBarVisibility(cview.ScrollBarAlways)
	p.searchResultTable.ShowFocus(true)
	p.searchResultTable.SetSelectionChangedFunc(p.updateDetails)
	p.searchResultSection.AddItem(p.searchResultTable, 0, 0, 1, 1, 0, 0, true)
	// p.searchResultSection.AddItem(p.searchResultTable, 0, 1, true)
	p.grid.AddItem(p.searchResultSection, 1, 0, 1, 1, 0, 0, true)

	detailsSection := cview.NewGrid()
	detailsSection.SetRows(-1)
	detailsSection.SetColumns(-1)
	detailsSection.SetTitle(" Details ")
	detailsSection.SetTitleAlign(cview.AlignLeft)
	detailsSection.SetBorder(true)

	p.descriptionView = cview.NewTextView()
	p.descriptionView.SetDynamicColors(true)
	p.descriptionView.SetRegions(true)
	p.descriptionView.SetWrap(true)
	p.descriptionView.SetWordWrap(true)
	detailsSection.AddItem(p.descriptionView, 0, 0, 1, 1, 0, 0, true)
	p.grid.AddItem(detailsSection, 2, 0, 1, 1, 0, 0, true)

	return p
}

func (p *searchPanel) readMessages() {
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
			go p.updateSearchResult(r)
		} else {
			m.DtoCastError()
		}

	default:
		m.UnsupportedTypeError()
	}
}

func (p *searchPanel) runSearch() {
	// p.searchResultSection.SetFocus()
	p.searchResult = make([]*dto.IAItem, 0)
	p.searchResultTable.Clear()
	// Disable Search Button here

	// Table Header
	r := 0
	headers := []string{
		"Title",
		"Files",
		"Duration",
		"Size",
	}
	for c, h := range headers {
		cell := cview.NewTableCell(h)
		textColor := tcell.ColorYellow.TrueColor()
		bgColor := tcell.ColorBlue.TrueColor()
		cell.SetTextColor(textColor)
		cell.SetBackgroundColor(bgColor)
		cell.SetAlign(cview.AlignLeft)
		p.searchResultTable.SetCell(r, c, cell)
	}
	p.searchResultTable.SetFixed(1, 0)
	p.searchResultTable.Select(1, 0)
	p.sendMessage(uiComponentName, controllerName, dto.SearchCommandType, dto.SearchCommand{SearchCondition: p.searchCriteria}, true)
}

func (p *searchPanel) clearSearchResults() {
	p.searchResult = make([]*dto.IAItem, 0)
	p.searchResultTable.Clear()
}

func (p *searchPanel) updateSearchResult(i *dto.IAItem) {
	logger.Debug(uiComponentName + ": Got AI Item: " + i.Title)
	p.searchResult = append(p.searchResult, i)
	r := p.searchResultTable.GetRowCount()
	cols := []string{
		i.Title,
		strconv.Itoa(i.FilesCount),
		strconv.Itoa(int(i.TotalLength)),
		strconv.Itoa(int(i.TotalSize)),
	}
	for c, col := range cols {
		p.searchResultTable.SetCell(r, c, cview.NewTableCell(col))
	}
	// p.searchResultTable.Select(r, 0)
	p.sendMessage(uiComponentName, "TUI", dto.CommandType, dto.Command{Command: "RedrawUI"}, true)
}

func (p *searchPanel) updateDetails(row int, col int) {
	if len(p.searchResult) > row {
		d := p.searchResult[row-1].Description
		p.descriptionView.SetText(d)
		p.descriptionView.ScrollToBeginning()
	}
}
