package controller

import (
	"time"

	"github.com/vpoluyaktov/audiobook_creator_IA/internal/mq"
)

type controller interface {
	readMessages()
}

type Dispatcher struct {
	dispatcher  *mq.Dispatcher
	controllers []controller
}

func NewDispatcher(dispatcher *mq.Dispatcher) *Dispatcher {
	c := &Dispatcher{}
	c.dispatcher = dispatcher
	c.controllers = append(c.controllers, NewSearchProcessor(c.dispatcher))
	return c
}

func (c *Dispatcher) startEventListener() {
	for {
		for _, p := range c.controllers {
			p.readMessages()
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (c *Dispatcher) Run() {
	go c.startEventListener()
}
