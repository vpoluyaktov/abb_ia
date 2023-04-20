package mq

import (
	"container/list"
	"sync"

	"github.com/vpoluyaktov/audiobook_creator_IA/internal/dto"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/logger"
)

type Dispatcher struct {
	mu         sync.Mutex
	recipients map[string]messageQueue
	listeners  map[string]CallBackFunc
}

type messageQueue struct {
	messages *list.List
}

type CallBackFunc func(*Message)

func NewDispatcher() *Dispatcher {
	d := &Dispatcher{}
	d.recipients = make(map[string]messageQueue)
	d.listeners = make(map[string]CallBackFunc)
	return d
}

func (d *Dispatcher) SendMessage(from string, to string, dtoType string, dto dto.Dto, async bool) {

	m := &Message{}
	m.From = from
	m.To = to
	m.Type = dtoType
	m.Dto = dto
	m.Async = async

	if async {
		// push message to queue
		if _, ok := d.recipients[m.To]; !ok {
			d.recipients[m.To] = messageQueue{list.New()}
		}
		d.mu.Lock()
		d.recipients[m.To].messages.PushBack(m)
		d.mu.Unlock()
		logger.Debug("MQ <-- async " + m.String())
	} else if _, ok := d.listeners[m.To]; ok {
		logger.Debug("MQ <-- sync  " + m.String())
		// call recepient method in blocking mode
		d.listeners[m.To](m)
		logger.Debug("MQ  sync --> " + m.String())
	}
}

func (d *Dispatcher) GetMessage(recipient string) *Message {
	var m *Message
	if _, ok := d.recipients[recipient]; ok {
		d.mu.Lock()
		e := d.recipients[recipient].messages.Front()
		if e != nil {
			d.recipients[recipient].messages.Remove(e)
			m = e.Value.(*Message)
			logger.Debug("MQ async --> " + m.String())
		}
		d.mu.Unlock()
	}
	return m
}

func (d *Dispatcher) RegisterListener(recipient string, callBackFunc CallBackFunc) {
	d.listeners[recipient] = callBackFunc
}
