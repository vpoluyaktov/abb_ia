package controller

import (
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/dto"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/event"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/logger"
)

const (
	componentName = "SearchProcessor"
)

type SearchProcessor struct {
	dispatcher *event.Dispatcher
}

func NewSearchProcessor(dispatcher *event.Dispatcher) *SearchProcessor {
	sp := &SearchProcessor{}
	sp.dispatcher = dispatcher
	sp.dispatcher.RegisterListener(componentName, sp.processMessage)
	return sp
}

func (p *SearchProcessor) readMessages() {
	m := p.dispatcher.GetMessage(componentName)
	if m != nil {
		p.processMessage(m)
	}
}

func (p *SearchProcessor) sendMessage(body interface{}) {
	m := &event.Message{}
	m.Sender = componentName
	m.Recipient = "SearchPanel"
	m.Type = "Text"
	m.Async = true
	m.Body = body
	p.dispatcher.SendMessage(m)
}

func (p *SearchProcessor) processMessage(m *event.Message) {
	if b, ok := m.Body.(dto.Button); ok {
		logger.Debug("SearchProcessor received a message from " + m.Sender + ": " + b.Name)
		p.sendMessage(b.Name + "Pressed")
	}
}
