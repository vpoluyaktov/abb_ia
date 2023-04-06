package mq

import (
	"container/list"
	"sync"
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

func (d *Dispatcher) SendMessage(message *Message) {
	if message.Async {
		// push message to queue
		if _, ok := d.recipients[message.To]; !ok {
			d.recipients[message.To] = messageQueue{list.New()}
		}
		d.mu.Lock()
		d.recipients[message.To].messages.PushFront(message)
		d.mu.Unlock()
	} else if _, ok := d.listeners[message.To]; ok {
		// call recepient method in blocking mode
		d.listeners[message.To](message)
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
		}
		d.mu.Unlock()
	}
	return m
}

func (d *Dispatcher) RegisterListener(recipient string, callBackFunc CallBackFunc) {
	d.listeners[recipient] = callBackFunc
}
