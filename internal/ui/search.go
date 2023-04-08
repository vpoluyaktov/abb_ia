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
	grid              *cview.Grid
	dispatcher        *mq.Dispatcher
	searchCriteria    string
	searchResult      []*dto.IAItem
	searchResultTable *cview.Table
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
	searchSection.AddItem(form, 0, 1, true)
	p.grid.AddItem(searchSection, 0, 0, 1, 1, 0, 0, true)

	searchResultSection := cview.NewFlex()
	searchResultSection.SetDirection(cview.FlexRow)
	searchResultSection.SetTitle(" Search result ")
	searchResultSection.SetTitleAlign(cview.AlignLeft)
	searchResultSection.SetBorder(true)
	p.searchResultTable = cview.NewTable()
	p.searchResultTable.SetSelectable(true, false)
	p.searchResultTable.SetScrollBarVisibility(cview.ScrollBarAlways)
	p.searchResultTable.ShowFocus(true)
	// p.searchResult.SetBorders(true)
	// color := tcell.ColorWhite.TrueColor()
	// cell := cview.NewTableCell("Cell value")
	// cell.SetTextColor(color)
	// cell.SetAlign(cview.AlignCenter)
	// p.searchResult.SetCell(0, 0, cell)
	searchResultSection.AddItem(p.searchResultTable, 0, 1, true)
	p.grid.AddItem(searchResultSection, 1, 0, 1, 1, 0, 0, true)
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
	p.searchResult = make([]*dto.IAItem, 0)
	p.searchResultTable.Clear()
	// Disable Search Button here
	p.sendMessage(uiComponentName, controllerName, dto.SearchCommandType, dto.SearchCommand{SearchCondition: p.searchCriteria}, true)
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
	p.sendMessage(uiComponentName, "TUI", dto.CommandType, dto.Command{Command: "RedrawUI"}, true)
}
