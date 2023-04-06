package controller

import (
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/dto"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/logger"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/mq"
)

const (
	controllerName  = "SearchController"
	uiComponentName = "SearchPanel"
)

type SearchController struct {
	dispatcher *mq.Dispatcher
}

func NewSearchProcessor(dispatcher *mq.Dispatcher) *SearchController {
	sp := &SearchController{}
	sp.dispatcher = dispatcher
	sp.dispatcher.RegisterListener(controllerName, sp.dispatchMessage)
	return sp
}

func (p *SearchController) readMessages() {
	m := p.dispatcher.GetMessage(controllerName)
	if m != nil {
		p.dispatchMessage(m)
	}
}

func (p *SearchController) sendMessage(from string, to string, dtoType string, dto dto.Dto, async bool) {
	m := &mq.Message{}
	m.From = from
	m.To = to
	m.Type = dtoType
	m.Dto = dto
	m.Async = async
	p.dispatcher.SendMessage(m)
}

func (p *SearchController) dispatchMessage(m *mq.Message) {
	switch t := m.Type; t {
	case dto.SearchCommandType:
		if c, ok := m.Dto.(dto.SearchCommand); ok {
			go p.performSearch(c)
		} else {
			m.DtoCastError()
		}

	default:
		m.UnsupportedTypeError()
	}
}

func (p *SearchController) performSearch(c dto.SearchCommand) {
	logger.Debug(controllerName + ": Received SearchCommand with condition: " + c.SearchCondition)
	// TODO: Run IA search
	p.sendMessage(controllerName, uiComponentName, dto.SearchResultType, dto.SearchResult{ItemName: "Some Item"}, true)
}
