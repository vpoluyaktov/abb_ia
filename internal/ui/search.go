package ui

import (
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
	grid       *cview.Grid
	dispatcher *mq.Dispatcher
}

func newSearchPanel(dispatcher *mq.Dispatcher) *searchPanel {
	p := &searchPanel{}
	p.dispatcher = dispatcher
	p.dispatcher.RegisterListener(uiComponentName, p.dispatchMessage)

	p.grid = cview.NewGrid()
	p.grid.SetRows(1, 1, 0, 3)
	p.grid.SetColumns(0)
	p.grid.AddItem(newText("SearchPanel"), 0, 0, 1, 1, 0, 0, true)
	p.grid.AddItem(newText("Body Search"), 1, 0, 1, 1, 0, 0, true)
	p.grid.AddItem(newButton("Search", p.runSearch), 3, 0, 1, 1, 0, 0, true)
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
	case dto.SearchResultType:
		if r, ok := m.Dto.(dto.SearchResult); ok {
			go p.updateSearchResult(r)
		} else {
			m.DtoCastError()
		}

	default:
		m.UnsupportedTypeError()
	}
}

func (p *searchPanel) runSearch() {
	// Disable Search Button here
	searchCondition := "NASA"
	p.sendMessage(uiComponentName, controllerName, dto.SearchCommandType, dto.SearchCommand{SearchCondition: searchCondition}, true)
}

func (p *searchPanel) updateSearchResult(r dto.SearchResult) {
	logger.Debug(uiComponentName + ": Got SearchResult: " + r.ItemName)

}
