package ui

import (
	"code.rocketnine.space/tslocum/cview"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/event"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/logger"
)

type searchPanel struct {
	grid       *cview.Grid
	dispatcher *event.Dispatcher
}

func newSearchPanel(dispatcher *event.Dispatcher) *searchPanel {
	p := &searchPanel{}
	p.dispatcher = dispatcher
	p.dispatcher.RegisterListener("SearchPanel", p.getMessage)

	p.grid = cview.NewGrid()
	p.grid.SetRows(1, 1, 0, 3)
	p.grid.SetColumns(0)
	p.grid.AddItem(newText("SearchPanel"), 0, 0, 1, 1, 0, 0, true)
	p.grid.AddItem(newText("Body Search"), 1, 0, 1, 1, 0, 0, true)
	p.grid.AddItem(newButton("NextPanel", func() {p.sendMessage("NextPanelButtonPressed")}), 3, 0, 1, 1, 0, 0, true)
	return p
}

func (p *searchPanel) getMessage(message *event.Message) {
	logger.Debug("SearchPanel received a message from " + message.Sender)
}

func (p *searchPanel) sendMessage(body interface{}) {
	m := &event.Message{}
	m.Sender = "SearchPanel"
	m.Recipient = "SearchController"
	m.Async = true
	m.Priority = 100
	m.Body = body
	p.dispatcher.SendMessage(m)
	
}