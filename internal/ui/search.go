package ui

import (
	"code.rocketnine.space/tslocum/cview"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/dto"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/event"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/logger"
)

const (
	componentName = "SearchPanel"
)

type searchPanel struct {
	grid       *cview.Grid
	dispatcher *event.Dispatcher
}

func newSearchPanel(dispatcher *event.Dispatcher) *searchPanel {
	p := &searchPanel{}
	p.dispatcher = dispatcher
	p.dispatcher.RegisterListener(componentName, p.processMessage)

	p.grid = cview.NewGrid()
	p.grid.SetRows(1, 1, 0, 3)
	p.grid.SetColumns(0)
	p.grid.AddItem(newText("SearchPanel"), 0, 0, 1, 1, 0, 0, true)
	p.grid.AddItem(newText("Body Search"), 1, 0, 1, 1, 0, 0, true)
	p.grid.AddItem(newButton("NextPanel", func() {p.sendMessage(dto.Button{Type: "Button", Name: "NextPanel", Event: "Pressed"})}), 3, 0, 1, 1, 0, 0, true)
	return p
}


func (p *searchPanel) readMessages() {
	m := p.dispatcher.GetMessage(componentName)
	if m != nil {
		p.processMessage(m)
	}
}

func (p *searchPanel) sendMessage(body interface{}) {
	m := &event.Message{}
	m.Sender = componentName
	m.Recipient = "SearchProcessor"
	m.Type = "Button"
	m.Async = true
	m.Body = body
	p.dispatcher.SendMessage(m)
	
}

func (p *searchPanel) processMessage(m *event.Message) {
	if s, ok := m.Body.(string); ok {
		logger.Debug("SearchPanel received a message from " + m.Sender + ": " + s)
	}
}