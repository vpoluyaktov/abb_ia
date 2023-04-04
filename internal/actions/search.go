package controller

import (
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/event"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/logger"
)

type SearchProcessor struct {
	dispatcher *event.Dispatcher
}

func NewSearchProcessor(dispatcher *event.Dispatcher) *SearchProcessor {
	sp := &SearchProcessor{}
	sp.dispatcher = dispatcher
	sp.dispatcher.RegisterListener("SearchController", sp.getMessage)
	return sp
}

func (p *SearchProcessor) ReadMessages() {
	m := p.dispatcher.GetMessage("SearchController")
	if m != nil {
		p.getMessage(m)
	}
}

func (p *SearchProcessor) getMessage(message *event.Message) {
	logger.Debug("SearchProcessor received a message from " + message.Sender)
	p.sendMessage("Pressed")
}

func (p *SearchProcessor) sendMessage(body interface{}) {
	m := &event.Message{}
	m.Sender = "SearchProcessor"
	m.Recipient = "SearchPanel"
	m.Async = false
	m.Priority = 100
	m.Body = body
	p.dispatcher.SendMessage(m)	
}