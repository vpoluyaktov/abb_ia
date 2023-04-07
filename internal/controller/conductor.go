package controller

import (
	"time"

	"github.com/vpoluyaktov/audiobook_creator_IA/internal/mq"
)

type controller interface {
	readMessages()
}

type Conductor struct {
	dispatcher  *mq.Dispatcher
	controllers []controller
}

func NewConductor(dispatcher *mq.Dispatcher) *Conductor {
	c := &Conductor{}
	c.dispatcher = dispatcher
	c.controllers = append(c.controllers, NewSearchController(c.dispatcher))
	return c
}

func (c *Conductor) startEventListener() {
	for {
		for _, p := range c.controllers {
			p.readMessages()
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (c *Conductor) Run() {
	go c.startEventListener()
}
